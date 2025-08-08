package generalSelectMenus

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/via-development/mr-poll/bot/internal/database"
	generalUtil "github.com/via-development/mr-poll/bot/internal/general-module/util"
)

func MrPollSelectMenu(interaction *events.ComponentInteractionCreate, db *database.GormDB) error {
	values := interaction.StringSelectMenuInteractionData().Values
	var embed discord.Embed
	switch values[0] {
	case "poll":
		embed = generalUtil.PollHelpPage()
	case "suggestion":
		embed = generalUtil.SuggestionHelpPage()
	default:
		embed = generalUtil.IntroductoryEmbed()
	}
	err := interaction.UpdateMessage(discord.MessageUpdate{
		Embeds: &[]discord.Embed{
			embed,
		},
	})
	return err
}
