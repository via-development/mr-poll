package pollUtil

import (
	"errors"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"mrpoll_bot/database"
	"slices"
)

var requiredPermissions = discord.PermissionSendMessages // | discord.PermissionViewChannel

// CreatePoll is used by both the website and poll command for creating polls, it checks permissions, creates poll message and saves to database.
func CreatePoll(client bot.Client, data database.PollData) (*discord.Message, error) {
	channelId := data.ChannelIdSnowflake()

	// Channel permission check
	channel, found := client.Caches().Channel(channelId)
	if !found {
		return nil, errors.New("channel is not in cache")
	}
	member, found := client.Caches().Member(*data.GuildIdSnowflake(), client.ApplicationID())
	if !found {
		m, err := client.Rest().GetMember(*data.GuildIdSnowflake(), client.ApplicationID())
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
	message, err := client.Rest().CreateMessage(channelId, discord.MessageCreate{
		Embeds:     MakePollEmbeds(data),
		Components: MakePollComponents(data, client),
	})
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

func CreatePollDM(interaction *events.ApplicationCommandInteractionCreate, data database.PollData) (*discord.Message, error) {
	// Send poll message in DM
	err := interaction.CreateMessage(discord.MessageCreate{
		Embeds:     MakePollEmbeds(data),
		Components: MakePollComponents(data, interaction.Client()),
	})
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

func CreatePollData(data database.PollData) error {
	return database.DB.Create(&data).Error
}

func VotePoll(pollData database.PollData, userId string, optionIds []int) (string, error) {
	action := ""
	for i, option := range pollData.Options {
		index := slices.Index(option.Voters, userId) // Find index of user's id in option's voters array

		if index != -1 { // They have already voted on the option

			if action == "added" {
				action = "swapped"
			} else {
				action = "removed"
			}

			if pollData.Type == database.SubmitChoiceType && len(option.Voters) <= 1 {
				database.DB.Delete(&option)
				pollData.Options = append(pollData.Options[:i], pollData.Options[i+1:]...)
			} else {
				pollData.Options[i].Voters = append(option.Voters[:index], option.Voters[index+1:]...)
				if err := database.DB.Save(&pollData.Options[i]).Error; err != nil {
					return "", err
				}
			}

		} else if slices.Index(optionIds, int(option.OptionId)) != -1 { // They haven't voted on the option yet

			if action == "removed" {
				action = "swapped"
			} else {
				action = "added"
			}

			pollData.Options[i].Voters = append(option.Voters, userId)
			if err := database.DB.Save(&pollData.Options[i]).Error; err != nil {
				return "", err
			}

		}
	}
	return action, nil
}

func UpdatePoll(interaction events.InteractionCreate) {
	//interaction.Client().Rest().UpdateInteractionResponse()
}
