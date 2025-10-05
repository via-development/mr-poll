package util

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
)

func CommandNotFoundEmbed(opts ...EmbedOption) discord.Embed {
	opts = append(opts, WithError("Command not found!"))
	return MakeEmbed(opts...)
}

func MakeSuccessEmbed(text string, opts ...EmbedOption) discord.Embed {
	opts = append(opts, WithSuccess(text), WithEmbedColor(GreenColor))
	return MakeEmbed(opts...)
}

func MakeErrorEmbed(text string, opts ...EmbedOption) discord.Embed {
	opts = append(opts, WithError(text), WithEmbedColor(RedColor))
	return MakeEmbed(opts...)
}

func MakeEmbed(opts ...EmbedOption) discord.Embed {
	embed := discord.Embed{
		Color: EmbedColor,
	}
	for _, opt := range opts {
		opt(&embed)
	}
	return embed
}

func WithEmbedColor(color int) func(*discord.Embed) {
	return func(embed *discord.Embed) {
		embed.Color = color
	}
}

func WithDescription(text string) func(*discord.Embed) {
	return func(embed *discord.Embed) {
		embed.Description = text
	}
}

//func WithPermissionsField() func(*discord.Embed) {
//	return func(embed *discord.Embed) {
//		embed.Fields = append(embed.Fields)
//	}
//}

func WithSuccess(text string) func(*discord.Embed) {
	return func(embed *discord.Embed) {
		embed.Description = fmt.Sprintf("<:e:%d> %s", CheckEmoji, text)
	}
}

func WithError(text string) func(*discord.Embed) {
	return func(embed *discord.Embed) {
		embed.Description = fmt.Sprintf("<:e:%d> %s", CrossEmoji, text)
	}
}

type EmbedOption func(*discord.Embed)
