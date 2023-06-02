package main

import (
	"fmt"
	"log"

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
	Session *dgo.Session
	Players map[string]*Playback
}

const (
	NO_PLAYER_AVAILABLE_ERR = "Some kind of error occurred, try again"
	NO_VOICE_CONNECTION_ERR = "No voice connection available"
	NO_TRACK_PLAYING_ERR    = "No track is being played"
	QUEUE_EMPTY_ERR         = "No track ready to be played"
	JOIN_CHANNEL_ERR        = "Some error occurred while trying to join channel, try again"
	BAD_COMMAND_ARG_ERR     = "Make sure to provide a valid command argument"
	SEEK_TOO_FAR_ERR        = "You went too far (the encoder might still be buffering)"
)

func NewClient(s *dgo.Session, guildIds []string) Client {
	c := Client{
		Players: make(map[string]*Playback),
		Session: s,
	}

	for _, gId := range guildIds {
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
