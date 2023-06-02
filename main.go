package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	dgo "github.com/bwmarrin/discordgo"
)

var RemoveCommands bool = false

const (
	ALIVE_COMMAND_NAME  = "alive"
	PLAY_COMMAND_NAME   = "play"
	NEXT_COMMAND_NAME   = "next"
	STOP_COMMAND_NAME   = "stop"
	PAUSE_COMMAND_NAME  = "pause"
	RESUME_COMMAND_NAME = "resume"
	SEEK_COMMAND_NAME   = "ff"
)

var commands = []*dgo.ApplicationCommand{
	{
		Name:        PLAY_COMMAND_NAME,
		Description: "Plays a song",
		Options: []*dgo.ApplicationCommandOption{
			{
				Name:        "input",
				Type:        dgo.ApplicationCommandOptionString,
				Description: "Raw media URL | YT web url | YT searchbar",
				Required:    true,
			},
		},
	},
	{
		Name:        SEEK_COMMAND_NAME,
		Description: "Fast forwards a song by a certain amount of seconds",
		Options: []*dgo.ApplicationCommandOption{
			{
				Name:        "input",
				Type:        dgo.ApplicationCommandOptionInteger,
				Description: "Amount of seconds to skip",
				Required:    true,
			},
		},
	},
	{
		Name:        STOP_COMMAND_NAME,
		Description: "Stops the current song",
	},
	{
		Name:        NEXT_COMMAND_NAME,
		Description: "Play the next song immediately",
	},
	{
		Name:        RESUME_COMMAND_NAME,
		Description: "Resumes a paused song",
	},
	{
		Name:        RESUME_COMMAND_NAME,
		Description: "Pauses a playing song",
	},
	{
		Name:        ALIVE_COMMAND_NAME,
		Description: "Am I alive? o.O",
	},
}

func main() {
	userHome, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	flags := flag.NewFlagSet("ndmb", flag.ContinueOnError)
	ffmpegPath := flags.String(
		"ffmpeg",
		"/usr/bin/ffmpeg",
		"Path to ffmpeg executable",
	)
	ytdlpPath := flags.String(
		"ytdlp",
		userHome+"/.local/bin/yt-dlp",
		"Path to ffmpeg executable",
	)
	token := flags.String(
		"token",
		"",
		"Discord bot token",
	)
	guildsStr := flags.String(
		"guilds",
		"",
		"A list of guild id for every discord server the bot will operate (comma separated)",
	)
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatal(err)
	}

	if *token == "" {
		fmt.Printf("Failed parsing bot token\n\n")
		if err := flags.Parse([]string{"-h"}); err != nil {
			log.Fatal(err)
		}
		os.Exit(1)
	}

	if *guildsStr == "" {
		fmt.Printf("Failed parsing bot token\n\n")
		if err := flags.Parse([]string{"-h"}); err != nil {
			log.Fatal(err)
		}
		os.Exit(1)
	}

	guilds := strings.Split(
		strings.ReplaceAll(*guildsStr, " ", ""),
		",",
	)

	if len(guilds) == 0 {
		fmt.Printf("Failed parsing bot guilds\n\n")
		if err := flags.Parse([]string{"-h"}); err != nil {
			log.Fatal(err)
		}
		os.Exit(1)
	}

	SetFfmpegPath(*ffmpegPath)
	SetYtdlpPath(*ytdlpPath)

	s, err := dgo.New("Bot " + *token)
	if err != nil {
		panic(err)
	}

	client := NewClient(s, guilds)

	s.AddHandler(func(s *dgo.Session, r *dgo.Ready) {
		username := s.State.User.Username + "#" + s.State.User.Discriminator
		log.Println("Logged in as: ", username)
	})

	s.AddHandler(func(s *dgo.Session, i *dgo.InteractionCreate) {
		commandName := i.ApplicationCommandData().Name
		log.Printf("User %s from channel %s invoked command: %s\n", i.Member.User.Username, i.GuildID, commandName)
		switch commandName {
		case ALIVE_COMMAND_NAME:
			client.AliveCommand(s, i)
		case PLAY_COMMAND_NAME:
			client.PlayCommand(s, i)
		case NEXT_COMMAND_NAME:
			client.NextCommand(s, i)
		case STOP_COMMAND_NAME:
			client.StopCommand(s, i)
		case PAUSE_COMMAND_NAME:
			client.PauseCommand(s, i)
		case RESUME_COMMAND_NAME:
			client.ResumeCommand(s, i)
		case SEEK_COMMAND_NAME:
			client.SeekCommand(s, i)
		default:
			log.Printf("%s no such command: %s\n", i.GuildID, commandName)
		}
	})

	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	for _, guildId := range guilds {
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
}
