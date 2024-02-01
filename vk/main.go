package vk

import (
	"log"
	"os"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

func Init() *longpoll.LongPoll {
	vk := api.NewVK(os.Getenv("VK_TOKEN"))

	lp, err := longpoll.NewLongPollCommunity(vk)
	if err != nil {
		log.Fatalf("LongPoll error: %s", err)
	} else {
		log.Printf("VK LongPoll started! Community ID: %d", lp.GroupID)
	}
	return lp
}
