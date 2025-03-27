package suggestionModule

import (
	suggestionCommands "mrpoll_bot/suggestion-module/commands"
	"mrpoll_bot/util"
)

var Module = &util.Module{
	Commands: map[string]util.ModuleCommand{
		"suggest": suggestionCommands.SuggestCommand,
	},
}
