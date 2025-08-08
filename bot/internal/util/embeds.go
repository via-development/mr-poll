package util

import (
	"github.com/disgoorg/disgo/discord"
)

func MakeSimpleEmbed(text string) discord.Embed {
	return discord.Embed{
		Description: text,
		Color:       EmbedColor,
	}
}

func MakeSuccessEmbed(text string) discord.Embed {
	return MakeSimpleEmbed("<:e:1268234822304792676> • " + text)
}
func CommandNotFoundEmbed() discord.Embed {
	return discord.Embed{
		Description: "Command not found!",
		Color:       EmbedColor,
	}
}
