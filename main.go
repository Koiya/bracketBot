package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/pelletier/go-toml"
	"os"
	"os/signal"
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
	// add a event handler
	dg.AddHandler(newMessage)
	if err := dg.Open(); err != nil {
		fmt.Println("Error opening connection,", err.Error())
		return
	}
	defer dg.Close()

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}
