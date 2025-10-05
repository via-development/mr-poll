package suggestionModule

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	"github.com/via-development/mr-poll/bot/internal/util"
)

func (m *SuggestionModule) SuggestionSubmitModal(interaction *events.ModalSubmitInteractionCreate) error {
	customId := interaction.Data.CustomID
	channelId := customId[len("suggest:submit:"):]

	var suggestionChannel schema.SuggestionChannel
	err := m.db.Find(&suggestionChannel, channelId).Error
	if err != nil {
		return err
	}

	suggestion := schema.Suggestion{
		GuildId:     interaction.GuildID().String(),
		ChannelId:   interaction.Channel().ID().String(),
		UserId:      interaction.User().ID.String(),
		Description: interaction.Data.Text("description"),
	}

	if title := interaction.Data.Text("title"); title != "" {
		suggestion.Title = &title
	}

	err = m.FetchSuggestionUser(&suggestion)
	if err != nil {
		return err
	}

	_, err = m.CreateSugggestion(&suggestion, &suggestionChannel)
	if err != nil {
		return nil
	}

	return interaction.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{
			util.MakeSuccessEmbed("The suggestion was created"),
		},
		Flags: discord.MessageFlagEphemeral,
	})
}
