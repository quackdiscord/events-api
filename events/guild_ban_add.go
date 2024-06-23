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
	eventType := "guild_ban_add"

	Events[eventType] = &Event{
		Name:    eventType,
		Handler: handleGuildBanAdd,
	}
}

type MemberBanned struct {
	Type    string        `json:"type"`
	GuildID string        `json:"guild_id"`
	User    *structs.User `json:"user"`
	Case    *structs.Case `json:"case"`
}

func handleGuildBanAdd(data string) {
	eventType := "guild_ban_add"

	// turn the string into a struct
	var msg MemberBanned
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

	if settings == nil || settings.MemberWebhookURL == "" {
		return
	}

	desc := fmt.Sprintf("**Member:** <@%s> (%s)", msg.User.ID, msg.User.Username)

	if msg.Case.ID != "" {
		desc += fmt.Sprintf("\n**Reason:** `%s`\n**Moderator:** <@%s> (%s)\n**Case ID:** %s", msg.Case.Reason, msg.Case.ModeratorID, msg.Case.ModeratorID, msg.Case.ID)
	}

	// make a post request to the webhook
	embed := structs.Embed{
		Title:       "Member banned",
		Color:       0xe75151,
		Description: desc,
		Author: structs.EmbedAuthor{
			Name: msg.User.Username,
			Icon: "https://cdn.discordapp.com/avatars/" + msg.User.ID + msg.User.Avatar + ".png",
		},
		Footer: structs.EmbedFooter{
			Text: fmt.Sprintf("User ID: %s", msg.User.ID),
		},
		Thumbnail: structs.EmbedThumbnail{
			URL: "https://cdn.discordapp.com/emojis/1064442673806704672.webp",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// check the length of the description
	if len(embed.Description) > 4096 {
		embed.Description = embed.Description[:4096]
	}

	err = lib.SendWHEmbed(settings.MemberWebhookURL, embed)
	if err != nil {
		log.WithError(err).Error("Error sending webhook message")
		return
	}

	log.Infof("Handled %s event", eventType)
}
