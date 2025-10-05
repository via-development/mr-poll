package poll_module

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	"github.com/via-development/mr-poll/bot/internal/util"
	"slices"
)

func (m *PollModule) PollOptionSubmitModal(interaction *events.ModalSubmitInteractionCreate) error {
	customId := interaction.Data.CustomID
	messageId := customId[len("poll:option-submit:"):]

	var pollData schema.Poll
	if err := m.db.Preload("Options").First(&pollData, messageId).Error; err != nil {
		_ = interaction.CreateMessage(PollNotFoundMessage())
		return err
	}

	err := m.FetchPollUser(&pollData)
	if err != nil {
		return err
	}

	userId := interaction.User().ID.String()
	optionName := interaction.Data.Text("answer")

	for i, option := range pollData.Options {
		if j := slices.Index(option.Voters, userId); j != -1 {
			if len(option.Voters) <= 1 {
				m.db.Delete(&option)
				pollData.Options = append(pollData.Options[:i], pollData.Options[i+1:]...)
			} else {
				option.Voters = append(option.Voters[:j], option.Voters[j+1:]...)

				m.db.Save(&option)
				pollData.Options[i] = option
			}
		}
	}

	optionId := uint(len(pollData.Options))
	for i := range len(pollData.Options) {
		if slices.IndexFunc(pollData.Options, func(op schema.PollOption) bool {
			return op.OptionId == uint(i)
		}) == -1 {
			optionId = uint(i)
			break
		}
	}

	optionData := schema.PollOption{
		Name:      optionName,
		OptionId:  optionId,
		Emoji:     util.Alpha[optionId],
		Voters:    []string{userId},
		MessageId: messageId,
		SubmitBy:  &userId,
	}
	m.db.Save(&optionData)
	pollData.Options = append(pollData.Options, optionData)

	pollEmbeds := m.MakePollEmbeds(&pollData)
	pollComponents := m.MakePollComponents(&pollData)
	return interaction.UpdateMessage(discord.MessageUpdate{
		Embeds:     &pollEmbeds,
		Components: &pollComponents,
	})
}
