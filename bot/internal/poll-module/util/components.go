package pollUtil

import (
	"fmt"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/via-development/mr-poll/bot/internal/database"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	"math"
)

var menuButton = discord.ButtonComponent{
	Emoji:    &discord.ComponentEmoji{ID: 1268234871034089543, Name: "e"},
	CustomID: "poll:menu",
	Style:    discord.ButtonStyleSecondary,
}

// MakePollComponents makes components for a poll with the poll data provided.
func MakePollComponents(client bot.Client, db *database.GormDB, data *schema.PollData) []discord.ContainerComponent {
	var components []discord.ContainerComponent
	switch data.Type {
	case schema.YesOrNoType, schema.SingleChoiceType:
		var options discord.ActionRowComponent
		for i, op := range data.Options {
			e := op.ApiEmoji()
			s := discord.ButtonStyleSecondary
			if data.Type == schema.YesOrNoType {
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
				Disabled: data.HasEnded,
			})
		}
		options = append(options, menuButton)
		for i := range (len(options) / 5) + 1 {
			upper := int(math.Min(float64(i*5+5), float64(len(options))))
			components = append(components, options[i*5:upper])
		}
	case schema.MultipleChoiceType:
		var options []discord.StringSelectMenuOption
		for _, opt := range data.Options {
			e := opt.ApiEmoji()
			desc := ""
			if data.CanSubmit && opt.SubmitBy != nil {
				user, err := db.FetchUser(client, *opt.SubmitBy)
				if err == nil && user != nil {
					desc = fmt.Sprint("Submitted by @", user.Username)
				}
			}
			options = append(options, discord.StringSelectMenuOption{
				Label:       opt.Name,
				Value:       fmt.Sprint("option-", opt.OptionId),
				Emoji:       &e,
				Description: desc,
			})
		}
		if len(options) < 10 && data.CanSubmit {
			options = append(options, discord.StringSelectMenuOption{
				Label: "Submit your answer!",
				Value: "option:submit",
			})
		}
		components = []discord.ContainerComponent{
			discord.ActionRowComponent{
				discord.StringSelectMenuComponent{
					CustomID:  "poll:opts",
					Options:   options,
					MaxValues: int(data.NumOfChoices),
					Disabled:  data.HasEnded,
				},
			},
			discord.ActionRowComponent{menuButton},
		}
	default:
		components[0] = discord.ActionRowComponent{
			discord.ButtonComponent{
				Label:    "Something went wrong!",
				CustomID: "oops",
				Style:    discord.ButtonStyleSecondary,
				Disabled: true,
			},
		}
	}

	return components
}
