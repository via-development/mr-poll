package generalModule

import (
	"github.com/disgoorg/disgo/events"
	generalCommands "mrpoll_bot/general-module/commands"
)

type GeneralModule struct {
	Commands map[string]func(*events.ApplicationCommandInteractionCreate) error
}

func makeGeneralModule() *GeneralModule {
	return &GeneralModule{
		Commands: map[string]func(*events.ApplicationCommandInteractionCreate) error{
			"help":    generalCommands.MrPollCommand,
			"mr-poll": generalCommands.MrPollCommand,
		},
	}
}

var Module = makeGeneralModule()
