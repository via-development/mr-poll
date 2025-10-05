package generalModule

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/via-development/mr-poll/bot/internal/util"
	"math/rand"
)

var titles = []string{
	"Uhh, whats the idea here?",
	"Well I don't see you suggesting a solution!",
	"Poll ended :) I can do it. Believe me please.",
	"Woah dude, whats your.. suggestion?!",
	"You can't steal a poll, that's not cool.",
	"BIG \\*\\*FLASHY\\*\\* TITLE",
	"You can take the man out of the poll but not the poll out of the man.",
	"These titles aren't funny.",
	"Unlimited polls for you, always",
	"you seemed to have created an awkward situation, i honestly never thought you would try to summon something on the level of a god as a servant. im afraid that's considered a foul system wise",
	"i think not therefore i am bot",
	"she don't know I'm thinking about another poll 🥶",
}

func IntroductoryEmbed() discord.Embed {
	title := titles[rand.Intn(len(titles))]
	return discord.Embed{
		Title: title,
		Image: &discord.EmbedResource{
			URL: "https://i.imgur.com/Vll0nQ4.png",
		},
		Color: util.EmbedColor,
	}
}

func PollHelpPage() discord.Embed {
	return discord.Embed{
		Title: "Poller",
		Color: util.EmbedColor,
	}
}

func SuggestionHelpPage() discord.Embed {
	return discord.Embed{
		Title: "Suggester",
		Color: util.EmbedColor,
	}
}

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
