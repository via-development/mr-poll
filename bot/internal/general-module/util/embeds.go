package generalUtil

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/via-development/mr-poll/bot/util"
)

func IntroductoryEmbed() discord.Embed {
	return discord.Embed{
		Title: "Welcome to Mr Poll!",
		Color: util.Config.EmbedColor,
	}
}

func PollHelpPage() discord.Embed {
	return discord.Embed{
		Title: "Poller",
		Color: util.Config.EmbedColor,
	}
}

func SuggestionHelpPage() discord.Embed {
	return discord.Embed{
		Title: "Suggestor",
		Color: util.Config.EmbedColor,
	}
}
