package suggestionModule

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	"github.com/via-development/mr-poll/bot/internal/util"
)

func (m *SuggestionModule) SuggestionsVoteButton(interaction *events.ComponentInteractionCreate) error {
	customId := interaction.Data.CustomID()
	op := customId[len("suggestions:"):]

	var suggestion schema.Suggestion
	err := m.db.Find(&suggestion, &schema.Suggestion{
		MessageId: interaction.Message.ID.String(),
	}).Error
	if err != nil {
		return err
	}

	switch op {
	case "voters":
		// TODO
	case "upvote", "downvote":
		action, err := m.VoteSuggestion(&suggestion, interaction.User().ID.String(), op == "upvote")
		if err != nil {
			return err
		}

		var suggestionChannel schema.SuggestionChannel
		err = m.db.Find(&suggestionChannel, &schema.SuggestionChannel{
			ChannelId: interaction.Channel().ID().String(),
		}).Error
		if err != nil {
			return err
		}

		err = interaction.UpdateMessage(discord.MessageUpdate{
			Embeds: &[]discord.Embed{
				m.MakeSuggestionEmbed(&suggestion, &suggestionChannel),
			},
		})

		if err != nil {
			return err
		}

		_, _ = interaction.Client().Rest.CreateFollowupMessage(interaction.Client().ID(), interaction.Token(), discord.MessageCreate{
			Flags: discord.MessageFlagEphemeral,
			Embeds: []discord.Embed{
				util.MakeSuccessEmbed(fmt.Sprintf("Your %s was %s", op, action)),
			},
		})

		return nil
	}
	return nil
}
