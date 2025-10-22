package generalModule

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func (m *GeneralModule) MrPollButton(interaction *events.ComponentInteractionCreate) error {
	page := interaction.Data.CustomID()[len("help:"):]
	var embed discord.Embed
	switch page {
	case "poll":
		embed = PollHelpPage()
	case "suggestion":
		embed = SuggestionHelpPage()
	default:
		embed = IntroductoryEmbed()
	}
	err := interaction.UpdateMessage(discord.MessageUpdate{
		Embeds: &[]discord.Embed{
			embed,
		},
		Components: &[]discord.LayoutComponent{
			HelpComponents(page == "back"),
		},
	})
	return err
}
