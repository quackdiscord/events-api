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
	eventType := "member_join"

	Events[eventType] = &Event{
		Name:    eventType,
		Handler: handleMemberJoin,
	}
}

type MemberJoin struct {
	Type    string          `json:"type"`
	GuildID string          `json:"guild_id"`
	Member  *structs.Member `json:"member"`
}

func handleMemberJoin(data string) {
	eventType := "member_join"

	// turn the string into a struct
	var msg MemberJoin
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

	desc := fmt.Sprintf("**Member:** <@%s> (%s)", msg.Member.User.ID, msg.Member.User.Username)

	// make a post request to the webhook
	embed := structs.Embed{
		Title:       "Member joined",
		Color:       0xeb459e,
		Description: desc,
		Author: structs.EmbedAuthor{
			Name: msg.Member.User.Username,
			Icon: "https://cdn.discordapp.com/avatars/" + msg.Member.User.ID + msg.Member.Avatar + ".png",
		},
		Footer: structs.EmbedFooter{
			Text: fmt.Sprintf("User ID: %s", msg.Member.User.ID),
		},
		Thumbnail: structs.EmbedThumbnail{
			URL: "https://cdn.discordapp.com/emojis/1064442704936828968.webp",
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
