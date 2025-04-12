package pollUtil

import (
	"github.com/disgoorg/disgo/discord"
)

func PollNotFoundMessage() discord.MessageCreate {
	return discord.MessageCreate{
		Flags:   discord.MessageFlagEphemeral,
		Content: "Cannot fetch this poll!",
	}
}

func PollOptionSubmitModel(messageId string) discord.ModalCreate {
	return discord.ModalCreate{
		Title:    "Submit your answer!",
		CustomID: "poll:option-submit:" + messageId,
		Components: []discord.ContainerComponent{
			discord.ActionRowComponent{
				discord.TextInputComponent{
					Label:    "Your Answer",
					CustomID: "answer",
					Style:    discord.TextInputStyleShort,
				},
			},
		},
	}
}
