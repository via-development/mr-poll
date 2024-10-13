package pollCommands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func PollCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	return interaction.CreateMessage(discord.MessageCreate{
		Content: "Hiya!",
		Flags:   discord.MessageFlagEphemeral,
	})
}
