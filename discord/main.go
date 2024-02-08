package discord

import (
	"log"
	"os"

	"github.com/CyberZoneDev/vk-ds/discord/commands"
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
		log.Fatalln("error opening connection,", err)
	}

	registerCommandHandlers(discord, commands.CommandHandlers)

	if err != nil {
		log.Fatalf("error creating slash commands: %v", err)
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

func registerCommandHandlers(s *discordgo.Session, commandHandlers []commands.CommandHandler) {
	log.Printf("registering %d command handlers", len(commandHandlers))
	for _, handler := range commandHandlers {
		cmd := handler.Command()
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			log.Fatalf("error creating command %s: %v", cmd.Name, err)
			continue
		}
	}
	s.AddHandler(func(dg *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			for _, handler := range commandHandlers {
				if handler.Command().Name == i.ApplicationCommandData().Name {
					handler.Handler(dg, i)
					return
				}
			}
		}
	})

}
