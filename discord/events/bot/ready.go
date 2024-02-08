package events

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func OnReady(d *discordgo.Session, event *discordgo.Ready) {
	d.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: "Киберспорт РТУ МИРЭА",
				Type: discordgo.ActivityTypeWatching,
				URL:  "https://vk.com/ssca_cybersport",
			},
		},
	})
	log.Printf("Discord bot %s is ready!", event.User.String())
}
