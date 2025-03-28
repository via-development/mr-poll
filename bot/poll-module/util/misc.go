package pollUtil

import (
	"github.com/disgoorg/disgo/discord"
	"mrpoll_bot/database"
)

// PollMessage generates the options for the poll's message, embed and component.
func PollMessage(data database.PollData) discord.MessageCreate {
	return discord.MessageCreate{
		Embeds:     MakePollEmbeds(data),
		Components: MakePollComponents(data),
	}
}

func PollNotFoundMessage() discord.MessageCreate {
	return discord.MessageCreate{
		Flags:   discord.MessageFlagEphemeral,
		Content: "Cannot fetch this poll!",
	}
}

func PollOptionSubmitModel() discord.ModalCreate {
	return discord.ModalCreate{
		Title:    "Submit your answer!",
		CustomID: "poll:submit",
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
