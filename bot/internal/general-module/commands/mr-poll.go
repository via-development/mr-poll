package generalCommands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/via-development/mr-poll/bot/internal/database"
	generalUtil "github.com/via-development/mr-poll/bot/internal/general-module/util"
)

func MrPollCommand(interaction *events.ApplicationCommandInteractionCreate, db *database.GormDB) error {
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
