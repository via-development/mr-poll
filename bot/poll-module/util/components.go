package pollUtil

import (
	"github.com/disgoorg/disgo/discord"
	pollDatabase "mrpoll_bot/poll-module/database"
)

var menuButton = discord.ButtonComponent{
	Emoji:    &discord.ComponentEmoji{ID: 1268234871034089543, Name: "e"},
	CustomID: "poll:menu",
	Style:    discord.ButtonStyleSecondary,
}

// MakePollComponents makes components for a poll with the poll data provided.
func MakePollComponents(data pollDatabase.PollData) discord.ActionRowComponent {
	switch data.Type {
	case pollDatabase.YesOrNoType:
		return discord.ActionRowComponent{
			discord.ButtonComponent{
				Emoji:    &discord.ComponentEmoji{ID: 1268234822304792676, Name: "e"},
				CustomID: "poll:option-0",
				Style:    discord.ButtonStyleSuccess,
			},
			discord.ButtonComponent{
				Emoji:    &discord.ComponentEmoji{ID: 1268234748988493905, Name: "e"},
				CustomID: "poll:option-1",
				Style:    discord.ButtonStyleDanger,
			},
			menuButton,
		}
	//case pollDatabase.SingleChoiceType:
	//	return discord.ActionRowComponent{
	//		discord.ButtonComponent{
	//			Label:    "a",
	//			CustomID: "a",
	//		},
	//	}
	//case pollDatabase.MultipleChoiceType:
	//	return discord.ActionRowComponent{
	//		discord.ButtonComponent{
	//			Label:    "a",
	//			CustomID: "a",
	//		},
	//	}
	//case pollDatabase.SubmitChoiceType:
	//	return discord.ActionRowComponent{
	//		discord.ButtonComponent{
	//			Label:    "a",
	//			CustomID: "a",
	//		},
	//	}
	default:
		return discord.ActionRowComponent{
			discord.ButtonComponent{
				Label:    "Something went wrong!",
				CustomID: "oops",
				Style:    discord.ButtonStyleSecondary,
				Disabled: true,
			},
		}
	}
}
