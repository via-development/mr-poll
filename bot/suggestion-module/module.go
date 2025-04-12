package suggestionModule

import (
	suggestionCommands "mrpoll_bot/suggestion-module/commands"
	"mrpoll_bot/util"
)

var Module = &util.Module{
	Name: "suggestion",
	Commands: map[string]util.ModuleCommand{
		"suggest": suggestionCommands.SuggestCommand,
	},
}
