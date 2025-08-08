package pollUtil

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
)

func MakePollText(pollData *schema.PollData) string {
	if pollData.HasEnded {
		if pollData.EnderUserId == nil {
			return fmt.Sprintf("⏱️ This poll was ended automatically. (<t:%d:R>)", pollData.EndAt.Unix())
		}
		return fmt.Sprintf("⏱️ This poll was ended by <@%s> (@%s). (<t:%d:R>)", *pollData.EnderUserId, pollData.EnderUser().Username, pollData.EndAt.Unix())
	}
	if pollData.EndAt != nil {
		if pollData.GuildId == nil {
			return "⏱️ This poll will not end automatically, this is a dm. How did you manage to do this?"
		}
		return fmt.Sprintf("⏱️ This poll will end automatically. (<t:%d:R>)", pollData.EndAt.Unix())
	}
	return ""
}

func PollNotFoundMessage() discord.MessageCreate {
	return discord.MessageCreate{
		Flags:   discord.MessageFlagEphemeral,
		Content: "Cannot fetch this poll!",
	}
}

func NotYourPollMessage() discord.MessageCreate {
	return discord.MessageCreate{
		Flags:   discord.MessageFlagEphemeral,
		Content: "This isn't your poll!",
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
