package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

type urls struct {
	Small string
}

type unsplashAPI struct {
	Id              string
	Alt_description string
	Urls            urls
}

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

	fmt.Println("Arke is now running!")
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

	if m.Content == "!serverinfo" {
		guildID := os.Getenv("GUILD_ID")
		guild, _ := s.Guild(guildID)

		r := "Server Info\n" +
			"Name: " + guild.Name + "\n" +
			"ID: " + guild.ID + "\n" +
			"Region: " + guild.Region + "\n" +
			"Description: " + guild.Description + "\n" +
			"Locale: " + guild.PreferredLocale + "\n"

		s.ChannelMessageSend(m.ChannelID, r)
	}

	if m.Content == "!botinfo" {
		hostname, _ := os.Hostname()

		r := "Bot Info\n" +
			"Name: " + s.State.User.Username + "\n" +
			"Discriminator: " + s.State.User.Discriminator + "\n" +
			"ID: " + s.State.User.ID + "\n" +
			"Host: " + hostname + "\n"

		s.ChannelMessageSend(m.ChannelID, r)
	}

	if m.Content == "!dog" || m.Content == "!cat" {
		t := trimFirstRune(m.Content)
		a := unsplashAPI{}

		url := "https://api.unsplash.com/photos/random/?query=" + t + "&client_id=S7DgKIlRASArvLuHj2hQ1tLQiisA1wzEYUEvg12FZsA"
		res, _ := http.Get(url)
		defer res.Body.Close()

		json.NewDecoder(res.Body).Decode(&a)

		r := a.Alt_description +
			"\n" +
			a.Urls.Small

		s.ChannelMessageSend(m.ChannelID, r)
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

	if sp[0] == "!whois" {
		guildID := os.Getenv("GUILD_ID")
		guildMembers, _ := s.GuildMembers(guildID, "", 10)

		if len(sp) == 1 {
			return
		}

		for _, member := range guildMembers {
			if member.User.Username == sp[1] {
				t, _ := member.JoinedAt.Parse()
				j := t.Local().Format(time.ANSIC)

				r := "User Info\n" +
					"Name: " + member.User.Username + "\n" +
					"Discriminator: " + member.User.Discriminator + "\n" +
					"ID: " + member.User.ID + "\n" +
					"Joined server: " + j + "\n" +
					"MFA status: " + strconv.FormatBool(member.User.MFAEnabled) + "\n" +
					"Verified status: " + strconv.FormatBool(member.User.Verified) + "\n"

				s.ChannelMessageSend(m.ChannelID, r)
			}
		}
	}

	if sp[0] == "!mute" || sp[0] == "!unmute" {
		guildID := os.Getenv("GUILD_ID")
		guildMembers, _ := s.GuildMembers(guildID, "", 10)

		if len(sp) == 1 {
			return
		}

		for _, member := range guildMembers {
			if member.User.Username == sp[1] {
				s.GuildMemberMute(guildID, member.User.ID, !member.Mute)

				if len(sp) == 3 {
					i, _ := strconv.Atoi(sp[2])
					ts := time.Duration(i)
					t := time.NewTimer(ts * time.Second)
					defer t.Stop()

					go func() {
						<-t.C
						s.GuildMemberMute(guildID, member.User.ID, !member.Mute)
					}()
				}

				return
			}
		}
	}

}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
