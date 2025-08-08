package util

import (
	"github.com/disgoorg/disgo/discord"
)

const embedColor = 0x40FFAC

func MakeSimpleEmbed(text string) discord.Embed {
	return discord.Embed{
		Description: text,
		Color:       embedColor,
	}
}

func MakeSuccessEmbed(text string) discord.Embed {
	return MakeSimpleEmbed("<:e:1268234822304792676> • " + text)
}
func CommandNotFoundEmbed() discord.Embed {
	return discord.Embed{
		Description: "Command not found!",
		Color:       embedColor,
	}
}
