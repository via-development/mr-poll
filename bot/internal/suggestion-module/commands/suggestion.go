package suggestionCommands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/via-development/mr-poll/bot/internal/database"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	suggestUtils "github.com/via-development/mr-poll/bot/internal/suggestion-module/utils"
	"github.com/via-development/mr-poll/bot/internal/util"
	"regexp"
	"slices"
	"strings"
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
		return suggestionChannelCommand(interaction, db)
	default:
		return interaction.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{
				util.CommandNotFoundEmbed(),
			},
		})
	}
}

var suggestionChannelNameRegex = regexp.MustCompile("^([A-z]|[0-9]|-|_)+$")

func suggestionChannelCommand(interaction *events.ApplicationCommandInteractionCreate, db *database.GormDB) error {
	cmdData := interaction.SlashCommandInteractionData()
	subcommand := *cmdData.SubCommandName

	name := cmdData.String("name")
	if match := suggestionChannelNameRegex.Match([]byte(name)); !match {
		return interaction.CreateMessage(discord.MessageCreate{
			Flags: discord.MessageFlagEphemeral,
			Embeds: []discord.Embed{
				util.MakeSimpleEmbed("Suggestion channel names can only include letters, numbers and dashes."),
			},
		})
	}
	name = strings.ToLower(name)
	guildId := interaction.GuildID().String()

	if subcommand == "add" {
		var suggestionChannels []schema.SuggestionChannel
		res := db.Find(&suggestionChannels, &schema.SuggestionChannel{GuildId: guildId})
		if res.Error != nil {
			return res.Error
		}

		// If more than 2 channels already exist, premium is required to add another
		if len(suggestionChannels) > 2 {
			return interaction.CreateMessage(discord.MessageCreate{
				Flags: discord.MessageFlagEphemeral,
				Embeds: []discord.Embed{
					util.MakeSimpleEmbed("ur too poor"),
				},
			})
		}

		channel := cmdData.Channel("channel")

		// Make sure the name isn't the same as a suggestion channel that already exists
		duplicateIn := slices.IndexFunc(suggestionChannels, func(sc schema.SuggestionChannel) bool {
			return sc.Name == name || sc.ChannelId == channel.ID.String()
		})

		if duplicateIn != -1 {
			return interaction.CreateMessage(discord.MessageCreate{
				Flags: discord.MessageFlagEphemeral,
				Embeds: []discord.Embed{
					util.MakeSimpleEmbed("This channel already exists"),
				},
			})
		}

		if !slices.Contains([]discord.ChannelType{discord.ChannelTypeGuildText, discord.ChannelTypeGuildNews}, channel.Type) {
			return interaction.CreateMessage(discord.MessageCreate{
				Flags: discord.MessageFlagEphemeral,
				Embeds: []discord.Embed{
					util.MakeSimpleEmbed("The channel must be a text channel!"),
				},
			})
		}

		suggestionChannel := &schema.SuggestionChannel{
			GuildId:   guildId,
			ChannelId: channel.ID.String(),
			Name:      name,
		}

		err := db.Create(suggestionChannel).Error
		if err != nil {
			return err
		}

		suggestionChannels = append(suggestionChannels, *suggestionChannel)
		err = suggestUtils.DeployGuildSuggestCommand(interaction.Client(), *interaction.GuildID(), suggestionChannels)
		if err != nil {
			return err
		}

		return interaction.CreateMessage(discord.MessageCreate{
			Flags: discord.MessageFlagEphemeral,
			Embeds: []discord.Embed{
				util.MakeSimpleEmbed("Suggestion channel added!"),
			},
		})
	} else if subcommand == "remove" {

	} else if subcommand == "edit" {

	}
	return nil
}
