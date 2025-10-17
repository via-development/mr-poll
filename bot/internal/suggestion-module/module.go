package suggestionModule

import (
	"github.com/via-development/mr-poll/bot/internal/config"
	"github.com/via-development/mr-poll/bot/internal/core"
	"github.com/via-development/mr-poll/bot/internal/database"
)

type SuggestionModule struct {
	core.Module

	db     *database.Database
	client *core.Client
	config *config.Config
}

func (m *SuggestionModule) Name() string {
	return "suggestion"
}

func (m *SuggestionModule) Commands() map[string]core.ModuleCommand {
	return map[string]core.ModuleCommand{
		"suggest":    m.SuggestCommand,
		"suggestion": m.SuggestionCommand,
	}
}

func (m *SuggestionModule) Buttons() []*core.ModuleComponent {
	return []*core.ModuleComponent{
		{"suggestions:", m.SuggestionsVoteButton},
	}
}

func (m *SuggestionModule) SelectMenus() []*core.ModuleComponent {
	return []*core.ModuleComponent{}
}

func (m *SuggestionModule) Modals() []*core.ModuleModal {
	return []*core.ModuleModal{
		{"suggest:submit:", m.SuggestionSubmitModal},
	}
}

func (m *SuggestionModule) MenuCommands() map[string]core.ModuleCommand {
	return map[string]core.ModuleCommand{}
}

func New(db *database.Database, client *core.Client, config *config.Config) *SuggestionModule {
	return &SuggestionModule{
		db:     db,
		client: client,
		config: config,
	}
}
