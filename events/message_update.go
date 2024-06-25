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
	eventType := "message_update"

	Events[eventType] = &Event{
		Name:    eventType,
		Handler: handleMessageUpdate,
	}
}

type MsgUpdate struct {
	Type           string                       `json:"type"`
	ID             string                       `json:"id"`
	Author         structs.LogUser              `json:"author"`
	GuildID        string                       `json:"guild_id"`
	ChannelID      string                       `json:"channel"`
	OldContent     string                       `json:"old_content"`
	NewContent     string                       `json:"new_content"`
	OldAttachments []*structs.MessageAttachment `json:"old_attachments"`
	NewAttachments []*structs.MessageAttachment `json:"new_attachments"`
}

func handleMessageUpdate(data string) {
	eventType := "message_update"

	// turn the string into a struct
	var msg MsgUpdate
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
	if len(msg.OldAttachments) == 0 && msg.OldContent == "" {
		return
	}

	desc := fmt.Sprintf("**Channel:** <#%s> (%s)\n**Author:** <@%s> (%s)", msg.ChannelID, msg.ChannelID, msg.Author.ID, msg.Author.Username)

	if msg.OldContent != msg.NewContent {
		desc += fmt.Sprintf("\n\n**Content:** ```diff\n- %s\n+%s```", msg.OldContent, msg.NewContent)
	}

	desc += fmt.Sprintf("\n[Jump to message](%s)", fmt.Sprintf("https://discord.com/channels/%s/%s/%s", msg.GuildID, msg.ChannelID, msg.ID))

	// make a post request to the webhook
	embed := structs.Embed{
		Title:       "Message edited",
		Color:       0x4ca99d,
		Description: desc,
		Author: structs.EmbedAuthor{
			Name: msg.Author.Username,
			Icon: msg.Author.AvatarURL,
		},
		Footer: structs.EmbedFooter{
			Text: fmt.Sprintf("Message ID: %s", msg.ID),
		},
		Thumbnail: structs.EmbedThumbnail{
			URL: "https://cdn.discordapp.com/emojis/1065110917962022922.webp",
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
