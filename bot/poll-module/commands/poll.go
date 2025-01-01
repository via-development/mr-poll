package pollCommands

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	baseUtil "mrpoll_bot/base-util"
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
	err := pollUtil.CreatePoll(interaction.Client(), pollUtil.PollCreateData{
		ChannelId: 979426067678888046,
		GuildId:   976147096757497937,
	})
	fmt.Println(err)
	return interaction.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{
			pollUtil.MakePollEmbed(),
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
