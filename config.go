package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

const DEFAULT_PATH = "./config/bot-data.json"

type BotConfig struct {
	Token    string   `json:"token"`
	GuildIds []string `json:"guildIds"`
}

func LoadConfig(optpath ...string) (out BotConfig) {
	path := DEFAULT_PATH
	if len(optpath) > 0 {
		path = optpath[0]
	}

	reader, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	buf, _ := ioutil.ReadAll(reader)
	json.Unmarshal(buf, &out)
	return out
}
