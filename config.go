package main

import (
	"encoding/json"
	"log"
	"os"
)

const DEFAULT_PATH = "./config/bot-data.json"

type BotConfig struct {
	Token    string   `json:"token"`
	GuildIds []string `json:"guildIds"`
}

// Loads the bot config data, panics on read errors
func LoadConfig(optpath ...string) (out BotConfig) {
	path := DEFAULT_PATH
	if len(optpath) > 0 {
		path = optpath[0]
	}

	reader, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(reader).Decode(&out)
	if err != nil {
		log.Fatal(err)
	}
	return out
}
