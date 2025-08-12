package generalButtons

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/via-development/mr-poll/bot/internal/database"
	generalUtil "github.com/via-development/mr-poll/bot/internal/general-module/util"
)

func MrPollButton(interaction *events.ComponentInteractionCreate, db *database.GormDB) error {
	page := interaction.Data.CustomID()[len("help:"):]
	var embed discord.Embed
	switch page {
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
		Components: &[]discord.ContainerComponent{
			generalUtil.HelpComponents(page == "back"),
		},
	})
	return err
}
