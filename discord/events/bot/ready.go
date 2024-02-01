package events

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func OnReady(d *discordgo.Session, event *discordgo.Ready) {
	d.UpdateWatchStatus(0, "https://vk.com/ssca_cybersport")
	log.Printf("Discord bot %s is ready!", event.User.String())
}
