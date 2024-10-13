package pollModule

import (
	"github.com/disgoorg/disgo/events"
	pollCommands "mrpoll_go/poll-module/commands"
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

const (
	FlagPollEnded               = 1 << iota
	FlagPollAnonymousVoting     = 1 << iota
	FlagPollDontShowCount       = 1 << iota
	FlagPollEmojiOptions        = 1 << iota
	FlagPollBonusVoteRolesStack = 1 << iota
)
