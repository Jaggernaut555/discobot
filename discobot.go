package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// variables used for command line parameters
var (
	Token     string
	ApiKey    string
	ChannelID string
)

// Pass discord bot token and youtube api key by flags -t and -k respectively
func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&ApiKey, "k", "", "Api key")
	flag.Parse()

	SetYoutubeKey(ApiKey)
}

func main() {

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// add a handler for when messages are posted
	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	defer dg.Close()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Do not let the bot talk to itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// specify the channel for the bot to use
	if m.Content == "!here" {
		ChannelID = m.ChannelID
	}

	if m.ChannelID != ChannelID {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "pong")
		return
	}

	if m.Content == "!help" {
		s.ChannelMessageSend(m.ChannelID, "use '!song [name]' to get a youtube link")
		return
	}

	if strings.HasPrefix(m.Content, "!song") {
		args := strings.SplitAfterN(m.Content, " ", 2)
		if len(args) < 2 {
			s.ChannelMessageSend(m.ChannelID, "Invalid request")
			return
		}
		result := YoutubeSearch(&args[1])
		if result == "" {
			s.ChannelMessageSend(m.ChannelID, "Could not find anything")
		} else {
			s.ChannelMessageSend(m.ChannelID, "http://youtu.be/"+result)
		}
	}
}
