package pollModals

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/via-development/mr-poll/bot/internal/database"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	pollUtil "github.com/via-development/mr-poll/bot/internal/poll-module/util"
	"github.com/via-development/mr-poll/bot/internal/util"
	"slices"
)

func PollOptionSubmitModal(interaction *events.ModalSubmitInteractionCreate, db *database.GormDB) error {
	customId := interaction.Data.CustomID
	messageId := customId[len("poll:option-submit:"):]

	var pollData schema.Poll
	if err := db.Preload("Options").First(&pollData, messageId).Error; err != nil {
		_ = interaction.CreateMessage(pollUtil.PollNotFoundMessage())
		return err
	}

	err := pollUtil.FetchPollUser(interaction.Client(), db, &pollData)
	if err != nil {
		return err
	}

	userId := interaction.User().ID.String()
	optionName := interaction.Data.Text("answer")

	for i, option := range pollData.Options {
		if j := slices.Index(option.Voters, userId); j != -1 {
			if len(option.Voters) <= 1 {
				db.Delete(&option)
				pollData.Options = append(pollData.Options[:i], pollData.Options[i+1:]...)
			} else {
				option.Voters = append(option.Voters[:j], option.Voters[j+1:]...)

				db.Save(&option)
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
	db.Save(&optionData)
	pollData.Options = append(pollData.Options, optionData)

	pollEmbeds := pollUtil.MakePollEmbeds(&pollData)
	pollComponents := pollUtil.MakePollComponents(interaction.Client(), db, &pollData)
	return interaction.UpdateMessage(discord.MessageUpdate{
		Embeds:     &pollEmbeds,
		Components: &pollComponents,
	})
}
