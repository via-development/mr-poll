package pollUtil

import (
	"errors"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"mrpoll_bot/database"
	pollDatabase "mrpoll_bot/poll-module/database"
)

var requiredPermissions = discord.PermissionSendMessages // | discord.PermissionViewChannel

// CreatePoll is used by both the website and poll command for creating polls, it checks permissions, creates poll message and saves to database.
func CreatePoll(client bot.Client, data pollDatabase.PollData) (*discord.Message, error) {
	channelId := data.ChannelIdSnowflake()

	// Channel permission check
	channel, found := client.Caches().Channel(channelId)
	if !found {
		return nil, errors.New("channel is not in cache")
	}
	member, found := client.Caches().Member(data.GuildIdSnowflake(), client.ApplicationID())
	if !found {
		m, err := client.Rest().GetMember(data.GuildIdSnowflake(), client.ApplicationID())
		if err != nil {
			return nil, err
		}
		member = *m
	}
	p := client.Caches().MemberPermissionsInChannel(channel, member)
	if p.Missing(requiredPermissions) {
		return nil, errors.New("i am missing permissions in the channel")
	}

	// Send poll message in channel
	message, err := client.Rest().CreateMessage(channelId, PollMessage(data))
	if err != nil {
		return nil, err
	}

	// Save poll to database
	data.MessageId = message.ID.String()
	if err = CreatePollData(data); err != nil {
		return nil, err
	}

	return message, nil
}

func CreatePollDM(interaction *events.ApplicationCommandInteractionCreate, data pollDatabase.PollData) (*discord.Message, error) {
	// Send poll message in DM
	err := interaction.CreateMessage(PollMessage(data))
	if err != nil {
		return nil, err
	}

	// Get interaction response
	message, err := interaction.Client().Rest().GetInteractionResponse(interaction.Client().ID(), interaction.Token())
	if err != nil {
		return nil, err
	}

	// Save poll to database
	data.MessageId = message.ID.String()
	if err = CreatePollData(data); err != nil {
		return nil, err
	}

	return message, nil
}

func CreatePollData(data pollDatabase.PollData) error {
	return database.DB.Create(&data).Error
}
