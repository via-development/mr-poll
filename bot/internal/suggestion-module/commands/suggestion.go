package suggestionCommands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/via-development/mr-poll/bot/internal/database"
	"github.com/via-development/mr-poll/bot/internal/util"
)

func SuggestionCommand(interaction *events.ApplicationCommandInteractionCreate, db *database.GormDB) error {
	subcommand := interaction.SlashCommandInteractionData().SubCommandGroupName
	if subcommand == nil {
		subcommand = interaction.SlashCommandInteractionData().SubCommandName
	}
	if subcommand == nil {
		return nil
	}

	switch *subcommand {
	case "channel":
		return suggestionChannelCommand(interaction)
	default:
		return interaction.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{
				util.CommandNotFoundEmbed(),
			},
		})
	}
}

func suggestionChannelCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	return nil
}
