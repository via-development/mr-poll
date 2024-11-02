package pollModule

import (
	"github.com/disgoorg/disgo/events"
	pollCommands "mrpoll_go/pkg/poll-module/commands"
)

type PollModule struct {
	Commands map[string]func(interaction *events.ApplicationCommandInteractionCreate) error
}

func makePollModule() *PollModule {
	return &PollModule{
		Commands: map[string]func(interaction *events.ApplicationCommandInteractionCreate) error{
			"poll": pollCommands.PollCommand,
		},
	}
}

var Module = makePollModule()
