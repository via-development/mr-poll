package suggestionModule

import (
	baseUtil "mrpoll_bot/base-util"
	suggestionCommands "mrpoll_bot/suggestion-module/commands"
)

var Module = &baseUtil.Module{
	Commands: map[string]baseUtil.ModuleCommand{
		"suggest": suggestionCommands.SuggestCommand,
	},
}
