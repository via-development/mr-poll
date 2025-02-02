package pollCommands

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	baseUtil "mrpoll_bot/base-util"
	pollDatabase "mrpoll_bot/poll-module/database"
	pollUtil "mrpoll_bot/poll-module/util"
)

func PollCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	subcommand := interaction.SlashCommandInteractionData().SubCommandName
	if subcommand == nil {
		return nil
	}

	switch *subcommand {
	case "yes-or-no", "multiple-choice", "single-choice", "submit-choice":
		return pollCreateCommand(interaction)
	case "list":
		return pollListCommand(interaction)
	case "end":
		return pollEndCommand(interaction)
	case "online":
		return pollOnlineCommand(interaction)
	default:
		return interaction.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{
				baseUtil.CommandNotFoundEmbed(),
			},
		})
	}
}

func pollCreateCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	pollData := pollDatabase.PollData{
		Type:      pollDatabase.YesOrNoType,
		ChannelId: interaction.Channel().ID().String(),
		Question:  "Holy Cow?",
		Options: []pollDatabase.PollOptionData{
			{
				OptionId: 0,
				Name:     "Yes",
				Voters:   make([]string, 0),
			},
			{
				OptionId: 1,
				Name:     "No",
				Voters:   make([]string, 0),
			},
		},
	}

	if guild, found := interaction.Guild(); found {
		pollData.GuildId = guild.ID.String()
	}

	if pollData.GuildId == "" {
		_, err := pollUtil.CreatePollDM(interaction, pollData)
		return err
	}

	message, err := pollUtil.CreatePoll(interaction.Client(), pollData)

	if err != nil {
		return err
	}

	return interaction.CreateMessage(discord.MessageCreate{
		Flags: discord.MessageFlagEphemeral,
		Embeds: []discord.Embed{
			baseUtil.MakeSimpleEmbed(fmt.Sprintf("Created [poll](https://discord.com/channels/%s/%s/%s)!", pollData.GuildId, message.ChannelID, message.ID)),
		},
	})
}

func pollOnlineCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	var guildId string
	if guild, found := interaction.Guild(); found {
		guildId = guild.ID.String()
	} else {
		guildId = interaction.Channel().ID().String()
	}

	return interaction.CreateMessage(discord.MessageCreate{
		Components: []discord.ContainerComponent{
			discord.ActionRowComponent{
				discord.ButtonComponent{ // Leaked?!?!
					URL:   "https://mrpoll.dev/polls/" + guildId,
					Style: discord.ButtonStyleLink,
					Label: "View Polls",
				},
			},
		},
	})
}

func pollEndCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	return nil
}

func pollListCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	return nil
}
