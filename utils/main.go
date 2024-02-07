package utils

import (
	"fmt"
	"strings"

	"github.com/SevereCloud/vksdk/v2/object"
	re "github.com/dlclark/regexp2"
)

func SortTags(text string, groupID int) string {
	var (
		mentions     = re.MustCompile(`\[([^\s\|]+)+\|(.*?)\]`, 0)
		subReg       = re.MustCompile(`^(?=id|club).*`, 0)
		hashtag      = re.MustCompile(`(?<![\S])\#(\S+(?<!\.))`, 0)
		hashtagGroup = re.MustCompile(`(?<!\[)\#(\S+(?<!\.))\@(\S+(?<!\.))`, 0)
		links        = re.MustCompile(`(?<!\(|\[|http:\/\/|https:\/\/|www\.)(?=(?<=\s)[A-Za-z0-9.]+\.[A-Za-z](?!\.))[^\s\)\]\!\?\;]+(?<!\.)`, 0)
	)

	oldText := text
	result, _ := mentions.FindStringMatch(text)
	var r0, r1, r2 string
	for result != nil {
		subResult, _ := subReg.FindStringMatch(result.GroupByNumber(1).String())
		if subResult != nil {
			r1 = "vk.com/" + subResult.String()
			r2 = result.GroupByNumber(2).String()
		} else {
			r1 = strings.Replace(result.GroupByNumber(1).String(), "https://", "", 1)
			r2 = result.GroupByNumber(2).String()
		}
		text = strings.Replace(oldText, result.String(), "["+r2+"](https://"+r1+")", 1)
		result, _ = mentions.FindStringMatch(text)
		oldText = text
	}

	result, _ = hashtagGroup.FindStringMatch(text)
	for result != nil {
		r0 = result.GroupByNumber(0).String()
		r1 = result.GroupByNumber(1).String()
		link := fmt.Sprintf("https://vk.com/wall%d?q=%%23%s", groupID, r1)
		text = strings.Replace(text, result.String(), "["+r0+"]("+link+")", 1)
		result, _ = hashtagGroup.FindStringMatch(text)
	}

	result, _ = hashtag.FindStringMatch(text)
	for result != nil {
		r0 = result.GroupByNumber(0).String()
		r1 = result.GroupByNumber(1).String()
		link := "https://vk.com/feed?section=search&q=%23" + r1
		text = strings.Replace(text, result.String(), "["+r0+"]("+link+")", 1)
		result, _ = hashtag.FindStringMatch(text)
	}

	result, _ = links.FindStringMatch(text)
	for result != nil {
		r0 = result.GroupByNumber(0).String()
		link := "https://" + r0
		text = strings.Replace(text, result.String(), "["+r0+"]("+link+")", 1)
		result, _ = links.FindStringMatch(text)
	}

	return text
}

func AttachmentsArray(attachments []object.WallWallpostAttachment) []string {
	var result []string
	for _, attachment := range attachments {
		var title string
		switch attachment.Type {
		case "doc":
			title = fmt.Sprintf("📄Файл: [%s](https://vk.com/doc%d_%d)", attachment.Doc.Title, attachment.Doc.OwnerID, attachment.Doc.ID)
			result = append(result, title)
		case "poll":
			title = fmt.Sprintf("📊Опрос: [%s](https://vk.com/poll%d_%d)", attachment.Poll.Question, attachment.Poll.OwnerID, attachment.Poll.ID)
			result = append(result, title)
		case "album":
			title = fmt.Sprintf("🖼️Альбом: [%s](https://vk.com/album%d_%d)", attachment.Album.Title, attachment.Album.OwnerID, attachment.Album.ID)
			result = append(result, title)
		case "video":
			if attachment.Video.Title == "Клип недоступен" || attachment.Video.Title == "Видео недоступно" {
				title = "🎥Видео"
			} else {
				title = fmt.Sprintf("🎥Видео: [%s](https://vk.com/video%d_%d)", attachment.Video.Title, attachment.Video.OwnerID, attachment.Video.ID)
			}
			result = append(result, title)
		}
	}
	return result
}
