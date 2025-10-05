package poll_module

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	"strconv"
)

func (m *PollModule) PollOptionSelectMenu(interaction *events.ComponentInteractionCreate) error {
	var poll schema.Poll
	if err := m.db.Preload("Options").First(&poll, interaction.Message.ID.String()).Error; err != nil {
		_ = interaction.CreateMessage(PollNotFoundMessage())
		return err
	}

	if poll.HasEnded {
		_ = interaction.CreateMessage(discord.MessageCreate{
			Flags:   discord.MessageFlagEphemeral,
			Content: "This poll has ended lol!",
		})
		return nil
	}

	err := m.FetchPollUser(&poll)
	if err != nil {
		return err
	}

	user := interaction.User()
	m.db.Save(&schema.User{
		UserId:      user.ID.String(),
		Username:    user.Username,
		DisplayName: user.GlobalName,
	})

	userId := interaction.User().ID.String()

	var selectedOptions []int
	{
		sv := interaction.StringSelectMenuInteractionData().Values

		for i := range sv {
			s := sv[i][len("option-"):]
			if s == "submit" {
				return interaction.Modal(PollOptionSubmitModel(interaction.Message.ID.String()))
			}
			n, _ := strconv.Atoi(s)
			selectedOptions = append(selectedOptions, n)
		}
	}

	action, err := m.VotePoll(&poll, userId, selectedOptions)
	if err != nil {
		return err
	}

	pollEmbeds := m.MakePollEmbeds(&poll)
	err = interaction.UpdateMessage(discord.MessageUpdate{
		Embeds: &pollEmbeds,
	})

	_ = interaction.CreateMessage(discord.MessageCreate{
		Flags:   discord.MessageFlagEphemeral,
		Content: fmt.Sprintf("Your vote was %s", action),
	})
	return err
}
