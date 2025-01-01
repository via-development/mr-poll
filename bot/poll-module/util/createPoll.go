package pollUtil

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
)

type PollOptionCreate struct {
	Id     int
	Name   string
	Emoji  string
	Voters []string
}

type PollCreateData struct {
	ChannelId uint64
	GuildId   uint64
	Question  string
	Options   []PollOptionCreate
}

var requiredPermissions = discord.PermissionSendMessages | discord.PermissionViewChannel

// CreatePoll is used by both the website and poll command for creating polls, it checks permissions, creates poll message and saves to database.
func CreatePoll(client bot.Client, data PollCreateData) error {
	// Channel Permission Check
	if data.GuildId != 0 {
		channel, found := client.Caches().Channel(snowflake.ID(data.ChannelId))
		if !found {
			return nil // TODO: ERROR
		}
		member, found := client.Caches().Member(snowflake.ID(data.GuildId), client.ApplicationID())
		if !found {
			m, err := client.Rest().GetMember(snowflake.ID(data.GuildId), client.ApplicationID())
			if err != nil {
				return err
			}
			member = *m
		}
		p := client.Caches().MemberPermissionsInChannel(channel, member)
		if !p.Has(requiredPermissions) {
			return nil // TODO: ERROR
		}
	}

	// Create poll embed
	// Create poll components
	// Send poll message in channel
	// Save poll to database
	return nil
}
