package util

import "github.com/disgoorg/disgo/discord"

func MakePollEmbed() discord.Embed {
	return discord.Embed{
		Title: "Poll",
	}
}
