package pollModals

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"mrpoll_bot/database"
	pollUtil "mrpoll_bot/poll-module/util"
	"mrpoll_bot/util"
	"slices"
)

func PollOptionSubmitModal(interaction *events.ModalSubmitInteractionCreate) error {
	customId := interaction.Data.CustomID
	messageId := customId[len("poll:option-submit:"):]

	var pollData database.PollData
	if err := database.DB.Preload("Options").First(&pollData, messageId).Error; err != nil {
		_ = interaction.CreateMessage(pollUtil.PollNotFoundMessage())
		return err
	}
	pollData.FetchUser(interaction.Client())

	userId := interaction.User().ID.String()
	optionName := interaction.Data.Text("answer")

	for i, option := range pollData.Options {
		if j := slices.Index(option.Voters, userId); j != -1 {
			if len(option.Voters) <= 1 {
				database.DB.Delete(&option)
				pollData.Options = append(pollData.Options[:i], pollData.Options[i+1:]...)
			} else {
				option.Voters = append(option.Voters[:j], option.Voters[j+1:]...)

				database.DB.Save(&option)
				pollData.Options[i] = option
			}
		}
	}

	optionId := uint(len(pollData.Options))
	for i := range len(pollData.Options) {
		if slices.IndexFunc(pollData.Options, func(op database.PollOptionData) bool {
			return op.OptionId == uint(i)
		}) == -1 {
			optionId = uint(i)
			break
		}
	}

	optionData := database.PollOptionData{
		Name:      optionName,
		OptionId:  optionId,
		Emoji:     util.Alpha[optionId],
		Voters:    []string{userId},
		MessageId: messageId,
		SubmitBy:  &userId,
	}
	database.DB.Save(&optionData)
	pollData.Options = append(pollData.Options, optionData)

	pollEmbeds := pollUtil.MakePollEmbeds(pollData)
	pollComponents := pollUtil.MakePollComponents(pollData, interaction.Client())
	return interaction.UpdateMessage(discord.MessageUpdate{
		Embeds:     &pollEmbeds,
		Components: &pollComponents,
	})
}
