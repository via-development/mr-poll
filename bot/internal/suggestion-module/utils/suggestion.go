package suggestUtils

import (
	"fmt"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
)

func DeployGuildSuggestCommand(client bot.Client, guildId snowflake.ID, suggestionChannels []schema.SuggestionChannel) error {
	gcommands, err := client.Rest().GetGuildCommands(client.ID(), guildId, true)
	if err != nil {
		return err
	}

	if len(suggestionChannels) == 0 {
		if len(gcommands) == 0 {
			return nil
		}

		var commandId *snowflake.ID
		for _, gcommand := range gcommands {
			if gcommand.Name() == "suggest" {
				ci := gcommand.ID()
				commandId = &ci
			}
		}

		if commandId == nil {
			return nil
		}

		return client.Rest().DeleteGuildCommand(client.ID(), guildId, *commandId)
	}

	var subcommands []discord.ApplicationCommandOption
	for _, sc := range suggestionChannels {
		subcommands = append(subcommands, discord.ApplicationCommandOptionSubCommand{
			Name:        sc.Name,
			Description: fmt.Sprintf("The %s suggestion channel", sc.Name),
		})
	}

	var commandId *snowflake.ID
	for _, gcommand := range gcommands {
		if gcommand.Name() == "suggest" {
			ci := gcommand.ID()
			commandId = &ci
		}
	}

	if commandId != nil {
		_, err = client.Rest().UpdateGuildCommand(client.ID(), guildId, *commandId, discord.SlashCommandUpdate{
			Options: &subcommands,
		})
	} else {
		_, err = client.Rest().CreateGuildCommand(client.ID(), guildId, discord.SlashCommandCreate{
			Name:        "suggest",
			Description: "The suggest command",
			Options:     subcommands,
		})
	}

	return err
}
