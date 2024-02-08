package commands

import (
	"github.com/CyberZoneDev/vk-ds/discord/commands/chat"
	"github.com/bwmarrin/discordgo"
)

type CommandHandler interface {
	Command() *discordgo.ApplicationCommand
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var CommandHandlers = []CommandHandler{
	&chat.SendPostHandler{},
}
