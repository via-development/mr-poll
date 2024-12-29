package generalCommands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"mrpoll_bot/general-module/util"
)

func MrPollCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	return interaction.CreateMessage(discord.MessageCreate{
		Embeds: []discord.Embed{
			generalUtil.IntroductoryEmbed(),
		},
		Components: []discord.ContainerComponent{
			discord.ActionRowComponent{
				generalUtil.HelpSelectMenu(),
			},
		},
	})
}
