package discord

import (
	"fmt"
	"log"
	"os"

	events "github.com/CyberZoneDev/vk-ds/discord/events/bot"
	"github.com/bwmarrin/discordgo"
)

type PostData struct {
	Author      string
	IconURL     string
	PostURL     string
	Text        string
	Picture     string
	Attachments []string
}

func Init() *discordgo.Session {
	// var (
	// 	GuildID = flag.String
	// )
	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Fatalf("Discord Init error: %s", err)
	}

	discord.AddHandlerOnce(events.OnReady)

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
	}

	return discord

}

func SendPost(d *discordgo.Session, data PostData) {
	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    data.Author,
			IconURL: data.IconURL,
			URL:     data.PostURL,
		},
		Description: data.Text,
		Image: &discordgo.MessageEmbedImage{
			URL: data.Picture,
		},
		Color: 0xf3660e,
	}

	_, err := d.ChannelMessageSendEmbed(os.Getenv("DISCORD_NEWS_CHANNEL"), embed)
	if err != nil {
		log.Printf("Failed to send embed message: %s", err)
	} else {
		log.Printf("Embed message sent | %s", data.PostURL)
	}
}
