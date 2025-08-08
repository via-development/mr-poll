package suggestionModule

import (
	suggestionCommands "github.com/via-development/mr-poll/bot/suggestion-module/commands"
	"github.com/via-development/mr-poll/bot/util"
)

var Module = &util.Module{
	Name: "suggestion",
	Commands: map[string]util.ModuleCommand{
		"suggest":    suggestionCommands.SuggestCommand,
		"suggestion": suggestionCommands.SuggestionCommand,
	},
}
