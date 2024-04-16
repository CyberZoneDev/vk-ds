package events

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func OnReady(d *discordgo.Session, event *discordgo.Ready) {
	err := d.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{
			{
				Name: "Киберспорт РТУ МИРЭА",
				Type: discordgo.ActivityTypeWatching,
				URL:  "https://vk.com/ssca_cybersport",
			},
		},
	})
	if err != nil {
		log.Fatalf("Error updating discord status: %s", err)
		return
	}
	log.Printf("Discord bot %s is ready!", event.User.String())
}
