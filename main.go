package main

import (
	"fmt"
	"io"
	"log"
	"ndmb/enc"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/Pauloo27/searchtube"
	dgo "github.com/bwmarrin/discordgo"
)

type PlaybackState int

const (
	PlaybackStateIdle PlaybackState = iota
	PlaybackStatePlaying
	PlaybackStateStopped
	PlaybackStatePaused
	PlaybackStateLooping
)

type Track struct {
	Title    string
	WebURL   string
	MediaURL string
}

type Playback struct {
	Track
	Queue             []Track
	PlayerCommandChan chan enc.Command
	state             PlaybackState
}

type Client struct {
	Session *dgo.Session
	Players map[string]Playback
}

func NewClient(s *dgo.Session, guildIds []string) Client {
	c := Client{
		Players: make(map[string]Playback),
		Session: s,
	}

	for _, gId := range guildIds {
		p := Playback{
			PlayerCommandChan: make(chan enc.Command),
			Queue:             make([]Track, 0),
		}
		c.Players[gId] = p
	}

	return c
}

var (
	s              *dgo.Session
	RemoveCommands bool = false
)

var (
	commands = []*dgo.ApplicationCommand{
		{
			Name:        "play",
			Description: "Plays a song",
			Options: []*dgo.ApplicationCommandOption{
				{
					Name:        "input",
					Type:        dgo.ApplicationCommandOptionString,
					Description: "Resource url or youtube search",
				},
			},
		},
		{
			Name:        "stop",
			Description: "Stops the current song",
		},
	}
)

func ResolveAudioSource(input string) (string, string, error) {
	webUrl := ""
	if IsYoutubeUrl(input) {
		webUrl = input
		mediaUrl, err := YoutubeMediaUrl(input)
		if err != nil {
			return "", webUrl, err
		}
		return mediaUrl, webUrl, nil
	}

	// Might be a direct http stream
	if strings.HasPrefix(input, "http") {
		return input, webUrl, nil
	}

	results, err := searchtube.Search(input, 1)
	if err != nil {
		return "", webUrl, err
	}

	webUrl = results[0].URL

	mediaUrl, err := YoutubeMediaUrl(webUrl)
	if err != nil {
		return "", webUrl, err
	}

	return mediaUrl, webUrl, nil
}

func FetchHttpMediaStream(mediaUrl string) (io.ReadCloser, int64, error) {
	req, err := http.NewRequest("GET", mediaUrl, nil)
	if err != nil {
		log.Println("Error preparing media stream request:", err)
		return nil, 0, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending media stream request:", err)
		return nil, 0, err
	}

	return resp.Body, resp.ContentLength, nil
}

func InteractionTextRespond(s *dgo.Session, i *dgo.InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Content: message,
		},
	})
}

func InteractionRespondDeferred(s *dgo.Session, i *dgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &dgo.InteractionResponse{
		Type: dgo.InteractionResponseDeferredChannelMessageWithSource,
		Data: &dgo.InteractionResponseData{
			Flags: dgo.MessageFlagsLoading,
		},
	})
}

func (c *Client) PlayCommand(s *dgo.Session, i *dgo.InteractionCreate) {
	genericErrorMessage := "Some error occurred while trying to join channel, try again."

	InteractionRespondDeferred(s, i)

	voiceChannelId := ""
	g, err := s.State.Guild(i.GuildID)
	if err != nil {
		InteractionTextRespond(s, i, "Couldn't any guild with id: "+i.GuildID) // Unlikely to happen
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
		if i.Interaction != nil {
			InteractionTextRespond(s, i, genericErrorMessage)
		}
		return
	}

	options := i.ApplicationCommandData().Options
	optionsMap := make(map[string]*dgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionsMap[opt.Name] = opt
	}

	userInput := optionsMap["input"].Value.(string) // assume a string, because yes
	mediaUrl, webUrl, err := ResolveAudioSource(userInput)
	if err != nil {
		InteractionTextRespond(s, i, genericErrorMessage)
		return
	}

	InteractionTextRespond(s, i, "Playing: "+webUrl)

	// Here buffering has to be implemented, also opus packets must be send at the correct rate
	ffmpegPath := "/usr/bin/ffmpeg"
	e := enc.NewEnc(enc.DefaultOptions(ffmpegPath))

	errCh := make(chan error)
	cmdCh := make(chan enc.Command)
	respCh := make(chan enc.Response)

	go func() {
		// cmdCh would ideally be saved in a guild queue
		e.GetOpusFrames(mediaUrl, enc.DefaultOptions(ffmpegPath), voiceConnection.OpusSend, errCh, cmdCh, respCh)
		close(errCh)
		close(respCh)
		close(cmdCh)
	}()

	go func() {
	loop:
		for {
			select {
			case err := <-errCh:
				log.Println("[ENCODER ERROR]:", err)
				break loop
			case response := <-respCh:
				log.Println("[RESPONSE]:", response)
			}
		}
	}()
}

func main() {
	botData := LoadConfig()
	s, err := dgo.New("Bot " + botData.Token)
	if err != nil {
		panic(err)
	}

	client := NewClient(s, botData.GuildIds)

	s.AddHandler(func(s *dgo.Session, r *dgo.Ready) {
		username := s.State.User.Username + "#" + s.State.User.Discriminator
		log.Println("Logged in as: ", username)
	})

	s.AddHandler(func(s *dgo.Session, i *dgo.InteractionCreate) {
		switch i.ApplicationCommandData().Name {
		case "play":
			client.PlayCommand(s, i)
		}
	})

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	for _, guildId := range botData.GuildIds {
		log.Println("Adding commands for guildId:", guildId)
		registeredCommands := make([]*dgo.ApplicationCommand, len(commands))
		for i, v := range commands {
			cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildId, v)
			log.Println("Created command", v.Name, "in", guildId)
			if err != nil {
				log.Panicf("Cannot create '%v' command: %v", v.Name, err)
			}
			registeredCommands[i] = cmd
		}
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Printf("Press Ctrl+C to exit")
	<-stop

	fmt.Println(botData.Token)
}
