package eventHandlers

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	generalModule "mrpoll_bot/general-module"
	pollModule "mrpoll_bot/poll-module"
)

func CommandHandler(e *events.ApplicationCommandInteractionCreate) {
	commandName := e.SlashCommandInteractionData().CommandName()

	command, ok := generalModule.Module.Commands[commandName]
	if ok {
		err := command(e)
		fmt.Println("Err: ", err)

		return
	}

	command, ok = pollModule.Module.Commands[commandName]
	if ok {
		err := command(e)
		fmt.Println("Err: ", err)

		return
	}

	_ = e.CreateMessage(discord.MessageCreate{
		Content: "I couldn't find that command!",
		Flags:   discord.MessageFlagEphemeral,
	})
}
