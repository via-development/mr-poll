package pollModule

import (
	"github.com/disgoorg/disgo/events"
	pollCommands "mrpoll_bot/poll-module/commands"
)

type PollModule struct {
	Commands map[string]func(*events.ApplicationCommandInteractionCreate) error
}

func makePollModule() *PollModule {
	return &PollModule{
		Commands: map[string]func(*events.ApplicationCommandInteractionCreate) error{
			"poll": pollCommands.PollCommand,
		},
	}
}

var Module = makePollModule()
