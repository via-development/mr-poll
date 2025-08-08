package suggestionModule

import (
	suggestionCommands "github.com/via-development/mr-poll/bot/internal/suggestion-module/commands"
	moduleUtil "github.com/via-development/mr-poll/bot/internal/util/module"
)

type SuggestionModule struct {
	moduleUtil.Module
}

func (m *SuggestionModule) Name() string {
	return "suggestion"
}

func (m *SuggestionModule) Commands() map[string]moduleUtil.ModuleCommand {
	return map[string]moduleUtil.ModuleCommand{
		"suggest":    suggestionCommands.SuggestCommand,
		"suggestion": suggestionCommands.SuggestionCommand,
	}
}

func (m *SuggestionModule) Buttons() []*moduleUtil.ModuleComponent {
	return []*moduleUtil.ModuleComponent{}
}

func (m *SuggestionModule) SelectMenus() []*moduleUtil.ModuleComponent {
	return []*moduleUtil.ModuleComponent{}
}

func (m *SuggestionModule) Modals() []*moduleUtil.ModuleModal {
	return []*moduleUtil.ModuleModal{}
}

func New() *SuggestionModule {
	return &SuggestionModule{}
}
