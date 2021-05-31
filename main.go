package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/niklasstich/AOCBot/handlers"
	"github.com/niklasstich/AOCBot/resources"
)

var session *discordgo.Session

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

	//closing shenanigans
	defer session.Close()
	lock := make(chan os.Signal)
	signal.Notify(lock, os.Interrupt, os.Kill)
	<-lock
	fmt.Println("shutting down gracefully")
}

func discordSetup(token string) error {
	var err error
	session, err = discordgo.New("Bot " + strings.Trim(token, "\n"))
	if err != nil {
		return err
	}
	err = session.Open()
	if err != nil {
		return err
	}
	session.AddHandler(handlers.CommandHandler)
	session.AddHandler(func(sess *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("Bot up!")
	})
	return nil
}
