package suggestionModule

import (
	"github.com/disgoorg/disgo/bot"
	baseUtil "mrpoll_bot/base-util"
	suggestionCommands "mrpoll_bot/suggestion-module/commands"
)

type SuggestionModule baseUtil.Module

var Module *SuggestionModule

func InitSuggestionModule(client *bot.Client) {
	Module = &SuggestionModule{
		Client: client,
		Commands: map[string]baseUtil.ModuleCommand{
			"suggest": suggestionCommands.SuggestCommand,
		},
		SelectMenus: []*baseUtil.ModuleSelectMenu{},
	}
}
