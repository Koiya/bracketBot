package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/pelletier/go-toml"
)

func main() {
	config, err := toml.LoadFile("config.toml")
	if err != nil {
		fmt.Println("Error ", err.Error())
		return
	}
	botToken := config.Get("Bot.token").(string)
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("Error creating Discord session,", err.Error())
		return
	}

	if err := dg.Open(); err != nil {
		fmt.Println("Error opening connection,", err.Error())
		return
	}
	defer dg.Close()
}
