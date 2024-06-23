package storage

import (
	"context"

	"github.com/quackdiscord/events-api/services"
	"github.com/quackdiscord/events-api/structs"
)

// find a log settings object by guild id in redis
func FindLogSettingsByID(id string) (*structs.LogSettings, error) {
	// get the hash table entry
	m, err := services.Redis.HGetAll(context.Background(), "ls_"+id).Result()
	if err != nil {
		return nil, err
	}

	// create a struct from the map
	g := structs.LogSettings{
		GuildID:           id,
		MessageChannelID:  m["message_channel_id"],
		MessageWebhookURL: m["message_webhook_url"],
		MemberChannelID:   m["member_channel_id"],
		MemberWebhookURL:  m["member_webhook_url"],
	}

	return &g, nil
}
