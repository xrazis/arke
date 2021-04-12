package main

import (
	"github.com/bwmarrin/discordgo"
	"os"
	"strconv"
	"strings"
	"time"
)

import . "marvin/api/requests"
import . "marvin/api/types"

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "!pong" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "Ping!")
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

		_, _ = s.ChannelMessageSend(m.ChannelID, r)

	}

	if m.Content == "!botinfo" {
		hostname, _ := os.Hostname()

		r := "Bot Info\n" +
			"Name: " + s.State.User.Username + "\n" +
			"Discriminator: " + s.State.User.Discriminator + "\n" +
			"ID: " + s.State.User.ID + "\n" +
			"Host: " + hostname + "\n"

		_, _ = s.ChannelMessageSend(m.ChannelID, r)
	}

	if m.Content == "!dog" || m.Content == "!cat" {
		at := UnsplashApi{}
		t := trimFirstRune(m.Content)
		unsplashAPIKey := os.Getenv("UNSPLASH_API_KEY")

		url := "https://api.unsplash.com/photos/random/?query=" + t + "&client_id=" + unsplashAPIKey

		ReadRespBody(url, &at)

		r := at.Alt_description +
			"\n" +
			at.Urls.Small

		_, _ = s.ChannelMessageSend(m.ChannelID, r)
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

		_ = s.ChannelMessagesBulkDelete(m.ChannelID, msgSlice)
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

				_, _ = s.ChannelMessageSend(m.ChannelID, r)
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
				_ = s.GuildMemberMute(guildID, member.User.ID, !member.Mute)

				if len(sp) == 3 {
					i, _ := strconv.Atoi(sp[2])
					ts := time.Duration(i)
					t := time.NewTimer(ts * time.Second)
					defer t.Stop()

					go func() {
						<-t.C
						_ = s.GuildMemberMute(guildID, member.User.ID, !member.Mute)
					}()
				}

				return
			}
		}
	}

}
