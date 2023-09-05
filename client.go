package main

import (
	"fmt"
	"log"
	"time"

	"ndmb/enc"

	dgo "github.com/bwmarrin/discordgo"
)

type Track struct {
	Title    string
	WebURL   string
	MediaURL string
}

type Playback struct {
	Track
	Player          *enc.Enc
	Queue           []Track
	CommandChannel  chan enc.Command
	ResponseChannel chan enc.Response
	ErrorChannel    chan error
	voiceConnection *dgo.VoiceConnection
}

type Client struct {
	Players        map[string]*Playback
	ActiveChannels map[string]string
}

const (
	NO_PLAYER_AVAILABLE_ERR = "Some kind of error occurred, try again"
	NO_VOICE_CONNECTION_ERR = "No voice connection available"
	NO_TRACK_PLAYING_ERR    = "No track is being played"
	QUEUE_EMPTY_ERR         = "No track ready to be played"
	JOIN_CHANNEL_ERR        = "Some error occurred while trying to join channel, try again"
	BAD_COMMAND_ARG_ERR     = "Make sure to provide a valid command argument"
	SEEK_TOO_FAR_ERR        = "You went too far (the encoder might still be buffering)"
	VOICE_IDLE_ERR          = "Failed disconnecting from idle channel connection"
	MAX_IDLE_SECONDS        = 300
)

func NewClient(s *dgo.Session, guildIds []string) Client {
	c := Client{
		Players:        make(map[string]*Playback),
		ActiveChannels: make(map[string]string, len(guildIds)),
	}

	for _, gId := range guildIds {
		c.ActiveChannels[gId] = ""

		p := &Playback{
			CommandChannel:  make(chan enc.Command),
			ResponseChannel: make(chan enc.Response),
			ErrorChannel:    make(chan error),
			Queue:           make([]Track, 0),
			Player:          enc.NewEnc(enc.DefaultOptions(GetFfmpegPath())),
		}

		// Whenever a track ends, play the next one
		p.Player.Listen(enc.PlayerEventTrackEnded, func(event enc.PlayerEvent) {
			if len(p.Queue) > 0 && p.voiceConnection != nil {
				nextTrack := p.Queue[0]
				p.Queue = p.Queue[1:]
				PlayMediaInVoiceChannel(nextTrack.MediaURL, p.Player, p.voiceConnection, p.ErrorChannel, p.CommandChannel, p.ResponseChannel)
			}
		})

		c.Players[gId] = p
	}

	return c
}

func (c *Client) ClientPlaybacksLogger(logInterval int) chan struct{} {
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-stop:
				return
			default:
			}

			for _, player := range c.Players {
				if player.voiceConnection == nil {
					continue
				}

				guildId := player.voiceConnection.GuildID
				select {
				case err := <-player.ErrorChannel:
					log.Printf("[PLAYER_ERR]: %v at guildId: %s\n", err, guildId)
				default:
				}
			}

			time.Sleep(time.Duration(int(time.Second) * logInterval))
		}
	}()

	return stop
}

// Logs to stdout the state of each player by guild id
func (c *Client) ClientLogger(each func() string, logInterval int) chan struct{} {
	stop := make(chan struct{})

	go func() {
		for {
			select {
			case <-stop:
				return
			default:
			}

			logValue := each()
			if logValue != "" {
				log.Println("[CLIENT_LOGGER]:", logValue)
			}
			time.Sleep(time.Duration(int(time.Second) * logInterval))
		}
	}()

	return stop
}

func (c *Client) StartDisconnectionTimmer(s *dgo.Session, tickEvery int) chan struct{} {
	stop := make(chan struct{})

	go func() {
		tickDuration := time.Duration(int64(time.Second) * int64(tickEvery))
		timers := make(map[string]int, len(c.Players))
		for guild := range c.Players {
			timers[guild] = 0
		}

		for {
			time.Sleep(tickDuration)

			select {
			case <-stop:
				for guild := range c.Players {
					if err := s.VoiceConnections[guild].Disconnect(); err != nil {
						log.Println("[VOICE_IDLE_ERR]:", VOICE_IDLE_ERR)
						continue
					}
				}
				return
			default:
			}

			for guild, player := range c.Players {
				voiceConnection := s.VoiceConnections[guild]
				if voiceConnection == nil {
					continue
				}

				members, err := s.GuildMembers(voiceConnection.GuildID, "", 0)
				shouldTick := player.Player.State == enc.PlayerStateIdle || (len(members) == 0 && err != nil)
				if err != nil {
					log.Println(err)
				}

				if shouldTick {
					timers[guild] += tickEvery
				} else {
					timers[guild] = 0
				}

				if timers[guild] >= MAX_IDLE_SECONDS {
					if err := s.VoiceConnections[guild].Disconnect(); err != nil {
						log.Println("[VOICE_IDLE_ERR]:", VOICE_IDLE_ERR)

						// do not reset timer, try again later
						continue
					}

					channelId := c.ActiveChannels[guild]
					_, err := s.ChannelMessageSend(channelId, "@here leaving channel as I've been idle for more than 5 minutes...")
					if err != nil {
						log.Printf("Failed sending disconnection message to channel %s with guildId %s\n", channelId, guild)
					}
					timers[guild] = 0
				}
			}
		}
	}()

	return stop
}

func (c *Client) PlayCommand(s *dgo.Session, i *dgo.InteractionCreate) {
	if err := InteractionRespondDeferred(s, i); err != nil {
		log.Printf(
			"Failed sending deferred response into guild: %s, error: %s",
			i.GuildID,
			err,
		)
	}

	voiceChannelId := ""
	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		err = InteractionTextUpdate(s, i, "Couldn't find any guild with id: "+i.GuildID) // Unlikely to happen
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				QUEUE_EMPTY_ERR,
				err,
			)
		}
		return
	}

	for _, vs := range g.VoiceStates {
		if vs.UserID == i.Member.User.ID {
			voiceChannelId = vs.ChannelID
			break
		}
	}

	voiceConnection, err := s.ChannelVoiceJoin(i.GuildID, voiceChannelId, false, true)
	if err != nil {
		if _, ok := s.VoiceConnections[i.GuildID]; ok {
			voiceConnection = s.VoiceConnections[i.GuildID]
		} else {
			err := InteractionTextUpdate(s, i, QUEUE_EMPTY_ERR)
			if err != nil {
				log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
					i.GuildID,
					QUEUE_EMPTY_ERR,
					err,
				)
			}
			return
		}
	}

	options := i.ApplicationCommandData().Options
	optionsMap := make(map[string]*dgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionsMap[opt.Name] = opt
	}

	userInput := optionsMap["input"].Value.(string) // assume a string, because yes
	track, err := ResolveAudioSource(userInput)
	if err != nil {
		err = InteractionTextUpdate(s, i, BAD_COMMAND_ARG_ERR)
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				BAD_COMMAND_ARG_ERR,
				err,
			)
		}
		return
	}

	var playback *Playback
	if p, ok := c.Players[i.GuildID]; ok {
		playback = p
	} else {
		err := InteractionTextUpdate(s, i, NO_PLAYER_AVAILABLE_ERR)
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				NO_PLAYER_AVAILABLE_ERR,
				err,
			)
		}
		return
	}

	if playback.Player.State == enc.PlayerStatePlaying || playback.Player.State == enc.PlayerStatePaused {
		msg := fmt.Sprintf("Track %s | %s added to queue", track.Title, track.WebURL)
		err := InteractionTextUpdate(s, i, msg)
		if err != nil {
			log.Printf("Failed sending client message to guild %s, client_message: %s, error: %s",
				i.GuildID,
				msg,
				err,
			)
		}
		playback.Queue = append(playback.Queue, track)
		return
	}

	playback.MediaURL = track.MediaURL
	playback.WebURL = track.WebURL
	playback.Title = track.Title

	msg := fmt.Sprintf("Now playing %s | %s", track.Title, track.WebURL)
	if err := InteractionTextUpdate(s, i, msg); err != nil {
		log.Printf("Failed sending client message to guild %s, client_message: %s, error: %s",
			i.GuildID,
			msg,
			err,
		)
	}

	playback.voiceConnection = voiceConnection
	PlayMediaInVoiceChannel(track.MediaURL, playback.Player,
		voiceConnection,
		playback.ErrorChannel,
		playback.CommandChannel,
		playback.ResponseChannel)
}

func (c *Client) NextCommand(s *dgo.Session, i *dgo.InteractionCreate) {
	if err := InteractionRespondDeferred(s, i); err != nil {
		log.Printf(
			"Failed sending deferred response into guild: %s, error: %s",
			i.GuildID,
			err,
		)
	}

	var playback *Playback
	if p, ok := c.Players[i.GuildID]; ok {
		playback = p
	} else {
		err := InteractionTextUpdate(s, i, NO_PLAYER_AVAILABLE_ERR)
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				NO_PLAYER_AVAILABLE_ERR,
				err,
			)

		}
		return
	}

	if playback.voiceConnection == nil {
		err := InteractionTextUpdate(s, i, NO_VOICE_CONNECTION_ERR)
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				NO_VOICE_CONNECTION_ERR,
				err,
			)
		}
		return
	}

	if playback.Player.State == enc.PlayerStateIdle || len(playback.Queue) == 0 {
		err := InteractionTextUpdate(s, i, QUEUE_EMPTY_ERR)
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				QUEUE_EMPTY_ERR,
				err,
			)
		}
		return
	}

	track := playback.Queue[0]
	playback.Queue = playback.Queue[1:]

	playback.MediaURL = track.MediaURL
	playback.WebURL = track.WebURL
	playback.Title = track.Title

	msg := fmt.Sprintf("Now playing %s | %s", track.Title, track.WebURL)
	err := InteractionTextUpdate(s, i, msg)
	if err != nil {
		log.Printf("Failed sending client message to guild %s, client_message: %s, error: %s",
			i.GuildID,
			msg,
			err,
		)
	}

	if playback.Player.State == enc.PlayerStatePlaying {
		playback.CommandChannel <- enc.CommandStop{}
	}

	PlayMediaInVoiceChannel(track.MediaURL, playback.Player,
		playback.voiceConnection,
		playback.ErrorChannel,
		playback.CommandChannel,
		playback.ResponseChannel)
}

func (c *Client) StopCommand(s *dgo.Session, i *dgo.InteractionCreate) {
	if err := InteractionRespondDeferred(s, i); err != nil {
		log.Printf(
			"Failed sending deferred response into guild: %s, error: %s",
			i.GuildID,
			err,
		)
	}

	var playback *Playback
	if p, ok := c.Players[i.GuildID]; ok {
		playback = p
	} else {
		err := InteractionTextUpdate(s, i, NO_PLAYER_AVAILABLE_ERR)
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				NO_PLAYER_AVAILABLE_ERR,
				err,
			)
		}
		return
	}

	if playback.voiceConnection == nil {
		err := InteractionTextUpdate(s, i, NO_VOICE_CONNECTION_ERR)
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				NO_VOICE_CONNECTION_ERR,
				err,
			)
		}
		return
	}

	if playback.Player.State != enc.PlayerStatePlaying && playback.Player.State != enc.PlayerStatePaused {
		err := InteractionTextUpdate(s, i, NO_TRACK_PLAYING_ERR)
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				NO_TRACK_PLAYING_ERR,
				err,
			)
		}
		return
	}

	playback.CommandChannel <- enc.CommandStop{}
	msg := fmt.Sprintf("Track %s | %s has been stopped", playback.Title, playback.WebURL)
	err := InteractionTextUpdate(s, i, msg)
	if err != nil {
		log.Printf("Failed sending client message to guild %s, client_message: %s, error: %s",
			i.GuildID,
			msg,
			err,
		)
	}
}

func (c *Client) PauseCommand(s *dgo.Session, i *dgo.InteractionCreate) {
	if err := InteractionRespondDeferred(s, i); err != nil {
		log.Printf(
			"Failed sending deferred response into guild: %s, error: %s",
			i.GuildID,
			err,
		)
	}

	var playback *Playback
	if p, ok := c.Players[i.GuildID]; ok {
		playback = p
	} else {
		err := InteractionTextUpdate(s, i, NO_PLAYER_AVAILABLE_ERR)
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				NO_PLAYER_AVAILABLE_ERR,
				err,
			)
		}
		return
	}

	if playback.voiceConnection == nil {
		err := InteractionTextUpdate(s, i, NO_VOICE_CONNECTION_ERR)
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				NO_VOICE_CONNECTION_ERR,
				err,
			)
		}
		return
	}

	if playback.Player.State != enc.PlayerStatePlaying && playback.Player.State != enc.PlayerStatePaused {
		err := InteractionTextUpdate(s, i, NO_TRACK_PLAYING_ERR)
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				NO_TRACK_PLAYING_ERR,
				err,
			)
		}
		return
	}

	playback.CommandChannel <- enc.CommandPause{}
	msg := fmt.Sprintf("Track %s | %s has been paused", playback.Title, playback.WebURL)
	err := InteractionTextUpdate(s, i, msg)
	if err != nil {
		log.Printf("Failed sending client message to guild %s, client_message: %s, error: %s",
			i.GuildID,
			msg,
			err,
		)
	}
}

func (c *Client) ResumeCommand(s *dgo.Session, i *dgo.InteractionCreate) {
	if err := InteractionRespondDeferred(s, i); err != nil {
		log.Printf(
			"Failed sending deferred response into guild: %s, error: %s",
			i.GuildID,
			err,
		)
	}

	var playback *Playback
	if p, ok := c.Players[i.GuildID]; ok {
		playback = p
	} else {
		err := InteractionTextUpdate(s, i, NO_PLAYER_AVAILABLE_ERR)
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				NO_PLAYER_AVAILABLE_ERR,
				err,
			)
		}
		return
	}

	if playback.voiceConnection == nil {
		err := InteractionTextUpdate(s, i, NO_VOICE_CONNECTION_ERR)
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				NO_VOICE_CONNECTION_ERR,
				err,
			)
		}
		return
	}

	if playback.Player.State != enc.PlayerStatePlaying && playback.Player.State != enc.PlayerStatePaused {
		err := InteractionTextUpdate(s, i, NO_TRACK_PLAYING_ERR)
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				NO_TRACK_PLAYING_ERR,
				err,
			)
		}
		return
	}

	playback.CommandChannel <- enc.CommandResume{}
	msg := fmt.Sprintf("Track %s | %s has been resumed", playback.Title, playback.WebURL)
	err := InteractionTextUpdate(s, i, msg)
	if err != nil {
		log.Printf("Failed sending client message to guild %s, client_message: %s, error: %s",
			i.GuildID,
			msg,
			err,
		)
	}
}

func (c *Client) SeekCommand(s *dgo.Session, i *dgo.InteractionCreate) {
	if err := InteractionRespondDeferred(s, i); err != nil {
		log.Printf(
			"Failed sending deferred response into guild: %s, error: %s",
			i.GuildID,
			err,
		)
	}

	options := i.ApplicationCommandData().Options
	optionsMap := make(map[string]*dgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionsMap[opt.Name] = opt
	}
	userInput := int(optionsMap["input"].Value.(float64))

	var playback *Playback
	if p, ok := c.Players[i.GuildID]; ok {
		playback = p
	} else {
		err := InteractionTextUpdate(s, i, NO_PLAYER_AVAILABLE_ERR)
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				NO_PLAYER_AVAILABLE_ERR,
				err,
			)
		}
		return
	}

	if playback.voiceConnection == nil {
		err := InteractionTextUpdate(s, i, NO_VOICE_CONNECTION_ERR)
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				NO_VOICE_CONNECTION_ERR,
				err,
			)
		}
		return
	}

	if playback.Player.State != enc.PlayerStatePlaying && playback.Player.State != enc.PlayerStatePaused {
		err := InteractionTextUpdate(s, i, NO_TRACK_PLAYING_ERR)
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				NO_TRACK_PLAYING_ERR,
				err,
			)
		}
		return
	}

	playback.CommandChannel <- enc.CommandGetDuration{}
	durationResponse := <-playback.ResponseChannel
	duration := int(durationResponse.(enc.ResponseDuration))

	playback.CommandChannel <- enc.CommandGetPlaybackTime{}
	currentTime := int((<-playback.ResponseChannel).(enc.ResponsePlaybackTime))

	cursor := currentTime + userInput
	if cursor > duration {
		err := InteractionTextUpdate(s, i, SEEK_TOO_FAR_ERR)
		if err != nil {
			log.Printf("Failed sending client error message to guild %s, client_error: %s, error: %s",
				i.GuildID,
				SEEK_TOO_FAR_ERR,
				err,
			)
		}
		return
	}

	minutes := cursor / 60
	seconds := cursor - (minutes * 60)

	playback.CommandChannel <- enc.CommandSeek(cursor)
	msg := fmt.Sprintf("Skipping track at %d:%d", minutes, seconds)
	err := InteractionTextUpdate(s, i, msg)
	if err != nil {
		log.Printf("Failed sending client message to guild %s, client_message: %s, error: %s",
			i.GuildID,
			msg,
			err,
		)
	}
}

func (c *Client) AliveCommand(s *dgo.Session, i *dgo.InteractionCreate) {
	msg := "I'm alive :)"
	err := InteractionTextRespond(s, i, msg)
	if err != nil {
		log.Printf("Failed sending client message to guild %s, client_message: %s, error: %s",
			i.GuildID,
			msg,
			err,
		)
	}
}

func (c *Client) LeaveCommand(s *dgo.Session, i *dgo.InteractionCreate) {
	connection := c.Players[i.GuildID].voiceConnection
	if connection == nil {
		ReportGenericError("No connection to end", s, i)
		return
	}

	// Stop the player/encoder if it's running for any reason
	if c.Players[i.GuildID].Player.State == enc.PlayerStatePaused ||
		c.Players[i.GuildID].Player.State == enc.PlayerStatePlaying {
		c.Players[i.GuildID].CommandChannel <- enc.CommandStop{}
	}

	err := connection.Disconnect()
	if err != nil {
		log.Println(err)
		ReportGenericError("No connection to end", s, i)
		return
	}
}

func ReportGenericError(msg string, session *dgo.Session, interaction *dgo.InteractionCreate) {
	err := InteractionTextRespond(session, interaction, msg)
	if err != nil {
		log.Printf("Failed sending client message to guild %s, client_message: %s, error: %s",
			interaction.GuildID,
			msg,
			err,
		)
	}
}
