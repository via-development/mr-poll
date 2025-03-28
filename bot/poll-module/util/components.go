package pollUtil

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"mrpoll_bot/database"
)

var menuButton = discord.ButtonComponent{
	Emoji:    &discord.ComponentEmoji{ID: 1268234871034089543, Name: "e"},
	CustomID: "poll:menu",
	Style:    discord.ButtonStyleSecondary,
}

// MakePollComponents makes components for a poll with the poll data provided.
func MakePollComponents(data database.PollData) discord.ActionRowComponent {
	// TODO: rewrite function
	switch data.Type {
	case database.YesOrNoType, database.SingleChoiceType:
		var options discord.ActionRowComponent
		for i, op := range data.Options {
			e := op.ApiEmoji()
			s := discord.ButtonStyleSecondary
			if data.Type == database.YesOrNoType {
				if i == 0 {
					s = discord.ButtonStyleSuccess
				} else {
					s = discord.ButtonStyleDanger
				}
			}
			options = append(options, discord.ButtonComponent{
				Emoji:    &e,
				CustomID: fmt.Sprint("poll:option-", i),
				Style:    s,
			})
		}
		options = append(options, menuButton)
		return options
	case database.MultipleChoiceType, database.SubmitChoiceType:
		var options []discord.StringSelectMenuOption
		for _, opt := range data.Options {
			e := opt.ApiEmoji()
			options = append(options, discord.StringSelectMenuOption{
				Label: opt.Name,
				Value: fmt.Sprint("option-", opt.OptionId),
				Emoji: &e,
			})
		}
		if len(options) < 10 && data.Type == database.SubmitChoiceType {
			options = append(options, discord.StringSelectMenuOption{
				Label: "Submit your answer!",
				Value: "option:submit",
			})
		}
		return discord.ActionRowComponent{
			discord.StringSelectMenuComponent{
				CustomID:  "poll:opts",
				Options:   options,
				MaxValues: int(data.NumOfChoices),
			},
		}
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
