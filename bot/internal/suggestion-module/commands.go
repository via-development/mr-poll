package suggestionModule

import (
	"errors"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
	"github.com/golittie/timeless"
	"github.com/via-development/mr-poll/bot/internal/database"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	"github.com/via-development/mr-poll/bot/internal/util"
	"gorm.io/gorm"
)

func (m *SuggestionModule) SuggestCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	commandData := interaction.SlashCommandInteractionData()
	if commandData.SubCommandName == nil {
		return nil
	}

	channelName := *commandData.SubCommandName
	var suggestionChannel schema.SuggestionChannel
	err := m.db.Find(&suggestionChannel, &schema.SuggestionChannel{GuildId: interaction.GuildID().String(), Name: channelName}).Error
	if err != nil {
		return err
	}

	return interaction.Modal(SuggestionSubmitModal(suggestionChannel.ChannelId))
}

const suggestionCommandPermissions = discord.PermissionManageChannels

func (m *SuggestionModule) SuggestionCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	if interaction.Member().Permissions.Missing(suggestionCommandPermissions) {
		return interaction.CreateMessage(discord.MessageCreate{
			Flags: discord.MessageFlagEphemeral,
			Embeds: []discord.Embed{
				util.MakeErrorEmbed("No perms :("),
			},
		})
	}

	subcommand := interaction.SlashCommandInteractionData().SubCommandGroupName
	if subcommand == nil {
		subcommand = interaction.SlashCommandInteractionData().SubCommandName
	}
	if subcommand == nil {
		return nil
	}

	switch *subcommand {
	case "channel":
		return m.suggestionChannelCommand(interaction)
	default:
		return interaction.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{
				util.CommandNotFoundEmbed(),
			},
		})
	}
}

var suggestionChannelNameRegex = regexp.MustCompile("^([A-z]|[0-9]|-|_)+$")

func (m *SuggestionModule) suggestionChannelCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	subcommand := *interaction.SlashCommandInteractionData().SubCommandName

	name := interaction.SlashCommandInteractionData().String("name")
	if match := suggestionChannelNameRegex.Match([]byte(name)); !match {
		return interaction.CreateMessage(discord.MessageCreate{
			Flags: discord.MessageFlagEphemeral,
			Embeds: []discord.Embed{
				util.MakeErrorEmbed("Suggestion channel names can only include letters, numbers and dashes."),
			},
		})
	}
	name = strings.ToLower(name)

	var suggestionChannels []schema.SuggestionChannel
	err := m.db.Find(&suggestionChannels, &schema.SuggestionChannel{GuildId: interaction.GuildID().String()}).Error
	if err != nil {
		return err
	}

	if subcommand == "add" {
		return m.suggestionChannelAddCommand(interaction, name, suggestionChannels)
	} else if subcommand == "remove" {
		return m.suggestionChannelRemoveCommand(interaction, name, suggestionChannels)
	} else if subcommand == "config" {
		return m.suggestionChannelConfigCommand(interaction, name, suggestionChannels)
	}

	return interaction.CreateMessage(discord.MessageCreate{
		Flags: discord.MessageFlagEphemeral,
		Embeds: []discord.Embed{
			util.CommandNotFoundEmbed(),
		},
	})
}

func (m *SuggestionModule) suggestionChannelAddCommand(interaction *events.ApplicationCommandInteractionCreate, name string, suggestionChannels []schema.SuggestionChannel) error {
	cmdData := interaction.SlashCommandInteractionData()
	// If more than 2 channels already exist, premium is required to add another
	if len(suggestionChannels) >= 2 {
		return interaction.CreateMessage(discord.MessageCreate{
			Flags: discord.MessageFlagEphemeral,
			Embeds: []discord.Embed{
				util.MakeErrorEmbed("ur too poor"),
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
				util.MakeErrorEmbed("This channel already exists"),
			},
		})
	}

	if !slices.Contains([]discord.ChannelType{discord.ChannelTypeGuildText, discord.ChannelTypeGuildNews}, channel.Type) {
		return interaction.CreateMessage(discord.MessageCreate{
			Flags: discord.MessageFlagEphemeral,
			Embeds: []discord.Embed{
				util.MakeErrorEmbed("The channel must be a text channel!"),
			},
		})
	}

	// Suggestion channel data creation

	suggestionChannel := schema.SuggestionChannel{
		GuildId:   interaction.GuildID().String(),
		ChannelId: channel.ID.String(),
		Name:      name,
		Panel: schema.SuggestionChannelPanel{
			Enabled: true,
			Embed:   false,
		},
	}

	if sendPanel, f := cmdData.OptBool("send-panel"); f {
		suggestionChannel.Panel.Enabled = sendPanel
	}

	if cooldown, f := cmdData.OptString("cooldown"); f {
		cooldownSecs := timeless.ParseTimeLength(cooldown, timeless.WithoutNegatives()).Seconds()
		if cooldownSecs < 30 || cooldownSecs > 24*60*60 {
			return interaction.CreateMessage(discord.MessageCreate{
				Flags: discord.MessageFlagEphemeral,
				Embeds: []discord.Embed{
					util.MakeErrorEmbed("Invalid time length provided for cooldown"),
				},
			})
		}
		suggestionChannel.Cooldown = int(cooldownSecs)
	}

	if embedColor, f := cmdData.OptString("embed-color"); f {
		ec, err := strconv.ParseInt(embedColor, 16, 16)
		if err != nil {
			return interaction.CreateMessage(discord.MessageCreate{
				Flags: discord.MessageFlagEphemeral,
				Embeds: []discord.Embed{
					util.MakeErrorEmbed("Invalid embed colour provided!"),
				},
			})
		}
		ec2 := int(ec)
		suggestionChannel.EmbedColor = &ec2
	}

	// Suggestion channel processing

	err := m.db.Create(&suggestionChannel).Error
	if err != nil {
		return err
	}

	suggestionChannels = append(suggestionChannels, suggestionChannel)
	err = m.DeployGuildSuggestCommand(*interaction.GuildID(), suggestionChannels)
	if err != nil {
		return err
	}

	if suggestionChannel.Panel.Enabled {
		mid, err := m.SendSuggestionChannelPanel(&suggestionChannel)
		if err != nil {
			return err
		}
		mid2 := mid.String()
		suggestionChannel.Panel.LastMessageId = &mid2
	}

	return interaction.CreateMessage(discord.MessageCreate{
		Flags: discord.MessageFlagEphemeral,
		Embeds: []discord.Embed{
			util.MakeSuccessEmbed("Suggestion channel added!"),
		},
	})
}

func (m *SuggestionModule) suggestionChannelRemoveCommand(interaction *events.ApplicationCommandInteractionCreate, name string, suggestionChannels []schema.SuggestionChannel) error {
	in := slices.IndexFunc(suggestionChannels, func(sc schema.SuggestionChannel) bool {
		return sc.Name == name
	})

	if in == -1 {
		return interaction.CreateMessage(discord.MessageCreate{
			Flags: discord.MessageFlagEphemeral,
			Embeds: []discord.Embed{
				util.MakeErrorEmbed("The suggestion channel doesn't exist!"),
			},
		})
	}

	err := m.db.Delete(&schema.SuggestionChannel{}, &schema.SuggestionChannel{
		ChannelId: suggestionChannels[in].ChannelId,
	}).Error
	if err != nil {
		return err
	}

	suggestionChannels = append(suggestionChannels[:in], suggestionChannels[in+1:]...)
	err = m.DeployGuildSuggestCommand(*interaction.GuildID(), suggestionChannels)
	if err != nil {
		return err
	}

	return interaction.CreateMessage(discord.MessageCreate{
		Flags: discord.MessageFlagEphemeral,
		Embeds: []discord.Embed{
			util.MakeSuccessEmbed("Suggestion channel deleted!"),
		},
	})
}

func (m *SuggestionModule) suggestionChannelConfigCommand(interaction *events.ApplicationCommandInteractionCreate, name string, suggestionChannels []schema.SuggestionChannel) error {
	in := slices.IndexFunc(suggestionChannels, func(sc schema.SuggestionChannel) bool {
		return sc.Name == name
	})

	if in == -1 {
		return interaction.CreateMessage(discord.MessageCreate{
			Flags: discord.MessageFlagEphemeral,
			Embeds: []discord.Embed{
				util.MakeErrorEmbed("The suggestion channel doesn't exist!"),
			},
		})
	}

	return nil
}

func (m *SuggestionModule) ApproveDenyCommand(interaction *events.ApplicationCommandInteractionCreate, db *database.Database) error {
	// TODO: perms check
	targetMessageId := interaction.SlashCommandInteractionData().String("message")

	var suggestion schema.Suggestion
	err := db.Find(&suggestion, &schema.Suggestion{
		MessageId: targetMessageId,
	}).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return interaction.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{
					util.MakeErrorEmbed("I couldn't find the suggestion!"),
				},
			})
		}
		return err
	}

	var suggestionChannel schema.SuggestionChannel
	err = db.Find(&suggestionChannel, &schema.SuggestionChannel{
		ChannelId: suggestion.ChannelId,
	}).Error
	if err != nil {
		return err
	}

	approved := interaction.Data.CommandName() == "approve"

	_, err = m.client.Rest().UpdateMessage(suggestion.ChannelIdSnowflake(), snowflake.MustParse(targetMessageId), discord.MessageUpdate{
		Embeds: &[]discord.Embed{
			m.MakeProcessedSuggestionEmbed(&suggestion, &suggestionChannel, approved),
		},
		Components: &[]discord.ContainerComponent{},
	})
	if err != nil {
		return interaction.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{
				util.MakeErrorEmbed("I couldn't update the suggestion!"),
			},
		})
	}

	state := "approved"
	if !approved {
		state = "denied"
	}

	return interaction.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{
			util.MakeSuccessEmbed("The suggestion was " + state + "."),
		},
	})
}
