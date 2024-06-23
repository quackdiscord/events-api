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
	eventType := "message_bulk_delete"

	Events[eventType] = &Event{
		Name:    eventType,
		Handler: handleMessageBulkDelete,
	}
}

type MsgBulkDelete struct {
	Type      string        `json:"type"`
	GuildID   string        `json:"guild_id"`
	ChannelID string        `json:"channel"`
	Messages  []BulkMessage `json:"messages"`
}

type BulkMessage struct {
	ID          string                       `json:"id"`
	Author      structs.LogUser              `json:"author"`
	Content     string                       `json:"content"`
	Attachments []*structs.MessageAttachment `json:"attachments"`
}

func handleMessageBulkDelete(data string) {
	eventType := "message_bulk_delete"

	// turn the string into a struct
	var msg MsgBulkDelete
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

	desc := fmt.Sprintf("**Channel:** <#%s> (%s)\n", msg.ChannelID, msg.ChannelID)

	for i, message := range msg.Messages {
		if message.Content != "" || len(message.Attachments) > 0 {
			desc += fmt.Sprintf("\n%d. <@%s> (%s)", i, message.Author.ID, message.Author.Username)

			if message.Content != "" {
				desc += fmt.Sprintf(" - `%s`", message.Content)
			}

			if len(message.Attachments) > 0 {
				desc += "\n> **Attachments:**"
				for _, attachment := range message.Attachments {
					desc += fmt.Sprintf("\n> - [%s](%s)", attachment.Filename, attachment.URL)
				}
			}
		}
	}

	// make a post request to the webhook
	embed := structs.Embed{
		Title:       fmt.Sprintf("%d messages deleted", len(msg.Messages)),
		Color:       0x373f69,
		Description: desc,
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
