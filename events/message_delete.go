package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/quackdiscord/events-api/lib"
	"github.com/quackdiscord/events-api/storage"
	"github.com/quackdiscord/events-api/structs"
	log "github.com/sirupsen/logrus"
)

func init() {
	eventType := "message_delete"

	Events[eventType] = &Event{
		Name:    eventType,
		Handler: handleMessageDelete,
	}
}

type MsgDelete struct {
	Type        string                       `json:"type"`
	ID          string                       `json:"id"`
	Author      structs.LogUser              `json:"author"`
	GuildID     string                       `json:"guild_id"`
	ChannelID   string                       `json:"channel"`
	Content     string                       `json:"content"`
	Attachments []*structs.MessageAttachment `json:"attachments"`
}

func handleMessageDelete(data string) {
	eventType := "message_delete"

	// turn the string into a struct
	var msg MsgDelete
	err := json.Unmarshal([]byte(data), &msg)
	if err != nil {
		log.WithError(err).Errorf("Error unmarshaling %s event", eventType)
		return
	}

	settings, err := storage.FindLogSettingsByID(msg.GuildID)
	if err != nil {
		log.WithError(err).Error("Error finding log settings")
		return
	}

	if settings == nil || settings.MessageWebhookURL == "" {
		return
	}

	// handle non attachment / non content messages
	if len(msg.Attachments) == 0 && msg.Content == "" {
		return
	}

	desc := fmt.Sprintf("**Channel:** <#%s> (%s)\n**Author:** <@%s> (%s)", msg.ChannelID, msg.ChannelID, msg.Author.ID, msg.Author.ID)

	if msg.Content != "" {
		desc += fmt.Sprintf("\n\n**Content:** `%s`", msg.Content)
	}

	if len(msg.Attachments) > 0 {
		desc += "\n\n**Attachments:**"
		for _, attachment := range msg.Attachments {
			desc += fmt.Sprintf("\n- [%s](%s)", attachment.Filename, attachment.URL)
		}
	}

	// make a post request to the webhook
	embed := structs.Embed{
		Title:       "Message deleted",
		Color:       0x914444,
		Description: desc,
		Author: structs.EmbedAuthor{
			Name: msg.Author.Username,
			Icon: msg.Author.AvatarURL,
		},
		Footer: structs.EmbedFooter{
			Text: fmt.Sprintf("Message ID: %s", msg.ID),
		},
		Thumbnail: structs.EmbedThumbnail{
			URL: "https://cdn.discordapp.com/emojis/1064444110334861373.webp",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// check the length of the description
	if len(embed.Description) > 4096 {
		embed.Description = embed.Description[:4096]
	}

	err = lib.SendWHEmbed(settings.MessageWebhookURL, embed)
	if err != nil {
		log.WithError(err).Error("Error sending webhook message")
		return
	}

	log.Infof("Handled %s event", eventType)
}
