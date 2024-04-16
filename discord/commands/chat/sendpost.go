package chat

import (
	"net/url"
	"os"
	"strings"

	"github.com/CyberZoneDev/vk-ds/utils"
	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/bwmarrin/discordgo"
)

type SendPostHandler struct{}

type PostData struct {
	Author      string
	IconURL     string
	PostURL     string
	Text        string
	Picture     string
	Attachments []string
}

var log = utils.NewLogger("chat")

func (h *SendPostHandler) Command() *discordgo.ApplicationCommand {
	permission := int64(discordgo.PermissionViewAuditLogs)
	return &discordgo.ApplicationCommand{
		Name:        "sendpost",
		Description: "Send post from VK to Discord",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "url",
				Description: "URL of the post",
				Required:    true,
			},
		},
		DefaultMemberPermissions: &permission,
	}
}

func (h *SendPostHandler) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	vk := api.NewVK(os.Getenv("VK_APP_TOKEN"))
	var data PostData
	data.PostURL = i.ApplicationCommandData().Options[0].StringValue()
	err := getPostData(&data, vk)
	if err != nil {
		log.Warnf("Error getting post data: %s", err)
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error getting post data",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Warnf("Error sending message: %s", err)
			return
		}
		return
	}
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
	_, err = s.ChannelMessageSendEmbed(i.ChannelID, embed)
	if err != nil {
		log.Warnf("Error sending discord message: %s", err)
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Done!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Warnf("Error sending discord message: %s", err)
		return
	}
}

func getPostData(data *PostData, vk *api.VK) error {
	u, err := url.Parse(data.PostURL)
	if err != nil {
		log.Warnf("URL parse error: %s", err)
		return err
	}
	path := u.Path
	if !strings.Contains(path, "wall") {
		log.Warnf("URL does not contain 'wall'")
		return err
	}
	parts := strings.Split(path, "wall")
	wallID := parts[1]

	post, err := vk.WallGetByID(api.Params{
		"posts": wallID,
	})
	if err != nil {
		log.Warnf("VK error: %s", err)
		return err
	}
	group, err := vk.GroupsGetByID(api.Params{"group_id": -post[0].OwnerID})
	if err != nil {
		log.Fatalf("Error getting group: %s", err)
	}

	text := utils.SortTags(post[0].Text, post[0].OwnerID)
	data.Author = group[0].Name
	data.IconURL = group[0].Photo200
	data.Text = text

	if len(post[0].Attachments) > 0 {
		if post[0].Attachments[0].Type == "photo" {
			data.Picture = post[0].Attachments[0].Photo.MaxSize().URL
		}
		attachments := utils.AttachmentsArray(post[0].Attachments)
		for attachmentsIndex, attachment := range attachments {
			if attachmentsIndex == 0 {
				data.Text += "\n\n"
			}
			data.Text += attachment + "\n"
		}
	}

	log.Infof("Post sended: %s", data.PostURL)
	return nil
}
