package generalUtil

import (
	"github.com/disgoorg/disgo/discord"
	baseUtil "mrpoll_bot/base-util"
)

func IntroductoryEmbed() discord.Embed {
	return discord.Embed{
		Title: "Welcome to Mr Poll!",
		Color: baseUtil.Config.EmbedColor,
	}
}

func PollPage() discord.Embed {
	return discord.Embed{
		Title: "Poller",
		Color: baseUtil.Config.EmbedColor,
	}
}

func SuggestionPage() discord.Embed {
	return discord.Embed{
		Title: "Suggestor",
		Color: baseUtil.Config.EmbedColor,
	}
}
