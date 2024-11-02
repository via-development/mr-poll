package baseUtil

import "github.com/disgoorg/disgo/discord"

func MakeSimpleEmbed(text string) discord.Embed {
	return discord.Embed{
		Description: text,
		Color:       Config.EmbedColor,
	}
}
