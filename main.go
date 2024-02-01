package main

import (
	"context"
	"strconv"

	// "fmt"
	"log"
	"os"
	"os/signal"

	"github.com/CyberZoneDev/vk-ds/discord"
	"github.com/CyberZoneDev/vk-ds/utils"
	"github.com/CyberZoneDev/vk-ds/vk"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/events"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil && os.Getenv("RUN_TYPE") != "CONTAINER" {
		log.Fatalf("Error loading .env: %s", err)
	}

	vkbot := vk.Init()
	discordbot := discord.Init()
	defer discordbot.Close()

	vkbot.WallPostNew(func(ctx context.Context, obj events.WallPostNewObject) {
		sortedtext := utils.SortTags(string(obj.Text))
		group, err := vkbot.VK.GroupsGetByID(api.Params{"group_id": -obj.OwnerID})
		if err != nil {
			log.Fatalf("Error getting group name: %s", err)
		}
		var data = discord.PostData{
			Author:  group[0].Name,
			IconURL: group[0].Photo200,
			PostURL: "https://vk.com/wall" + strconv.Itoa(obj.OwnerID) + "_" + strconv.Itoa(obj.ID),
			Text:    sortedtext,
		}
		if len(obj.Attachments) > 0 {
			if obj.Attachments[0].Type == "photo" {
				data.Picture = obj.Attachments[0].Photo.MaxSize().URL
			}
			attachments := utils.AttachmentsArray(obj.Attachments)
			for attachmentsIndex, attachment := range attachments {
				if attachmentsIndex == 0 {
					data.Text += "\n\n"
				}
				data.Text += attachment + "\n"
			}
		}
		discord.SendPost(discordbot, data)
	})
	if err := vkbot.Run(); err != nil {
		log.Fatal(err)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Graceful shutdown")
}
