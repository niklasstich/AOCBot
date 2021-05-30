package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/niklasstich/AOCBot/handlers"
	"github.com/niklasstich/AOCBot/resources"
)

func main() {
	fmt.Println("starting bot")
	config, err := resources.Config()
	if err != nil {
		panic(err)
	}
	err = discordSetup(config.BotToken)
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening for input..")
	lock := make(chan struct{})
	<-lock
}

func discordSetup(token string) error {
	discord, err := discordgo.New("Bot " + strings.Trim(token, "\n"))
	if err != nil {
		return err
	}
	err = discord.Open()
	if err != nil {
		return err
	}
	discord.AddHandler(handlers.CommandHandler)
	return nil
}
