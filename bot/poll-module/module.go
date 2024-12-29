package pollModule

import (
	"github.com/disgoorg/disgo/bot"
	baseUtil "mrpoll_bot/base-util"
	pollCommands "mrpoll_bot/poll-module/commands"
)

type PollModule baseUtil.Module

func (module *PollModule) SendPollMessage() {

}

var Module *PollModule

func InitPollModule(client *bot.Client) {
	Module = &PollModule{
		Client: client,
		Commands: map[string]baseUtil.ModuleCommand{
			"poll": pollCommands.PollCommand,
		},
	}
}
