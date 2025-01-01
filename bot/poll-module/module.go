package pollModule

import (
	baseUtil "mrpoll_bot/base-util"
	pollCommands "mrpoll_bot/poll-module/commands"
)

var Module = &baseUtil.Module{
	Commands: map[string]baseUtil.ModuleCommand{
		"poll": pollCommands.PollCommand,
	},
}
