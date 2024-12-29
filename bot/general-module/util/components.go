package generalUtil

import "github.com/disgoorg/disgo/discord"

func HelpSelectMenu() *discord.StringSelectMenuComponent {
	return &discord.StringSelectMenuComponent{
		CustomID:    "help:",
		Placeholder: "Select a module page",
		Options: []discord.StringSelectMenuOption{
			{Label: "Home", Value: "home"},
			{Label: "Polling", Value: "poll"},
			{Label: "Suggestion", Value: "suggestion"},
		},
	}
}
