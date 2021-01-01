package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("DISCORD_TOKEN")

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "!pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

	sp := strings.Split(m.Content, " ")

	if sp[0] == "!delete" {
		var msgSlice []string
		var r int

		if len(sp) > 1 {
			r, _ = strconv.Atoi(sp[1])
		} else {
			r = 2
		}

		msg, _ := s.ChannelMessages(m.ChannelID, r, "", "", "")

		for _, message := range msg {
			msgSlice = append(msgSlice, message.ID)
		}

		s.ChannelMessagesBulkDelete(m.ChannelID, msgSlice)
	}

}
