package generalUtil

import "github.com/disgoorg/disgo/discord"

func HelpComponents(isHome bool) discord.ActionRowComponent {
	if !isHome {
		return discord.NewActionRow(
			discord.ButtonComponent{
				CustomID: "help:back",
				Label:    "⇽ Back",
				Style:    discord.ButtonStyleDanger,
			},
		)
	}
	return discord.NewActionRow(
		discord.ButtonComponent{
			CustomID: "help:poll",
			Label:    "Polls →",
			Style:    discord.ButtonStyleSecondary,
		},
		discord.ButtonComponent{
			CustomID: "help:suggestion",
			Label:    "Suggestions →",
			Style:    discord.ButtonStyleSecondary,
		},
	)
}
