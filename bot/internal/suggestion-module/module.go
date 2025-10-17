package suggestionModule

import (
	"github.com/via-development/mr-poll/bot/internal"
	"github.com/via-development/mr-poll/bot/internal/config"
	"github.com/via-development/mr-poll/bot/internal/database"
)

type SuggestionModule struct {
	internal.Module

	db     *database.Database
	client *internal.MPBot
	config *config.Config
}

func (m *SuggestionModule) Name() string {
	return "suggestion"
}

func (m *SuggestionModule) Commands() map[string]internal.ModuleCommand {
	return map[string]internal.ModuleCommand{
		"suggest":    m.SuggestCommand,
		"suggestion": m.SuggestionCommand,
	}
}

func (m *SuggestionModule) Buttons() []*internal.ModuleComponent {
	return []*internal.ModuleComponent{
		{"suggestions:", m.SuggestionsVoteButton},
	}
}

func (m *SuggestionModule) SelectMenus() []*internal.ModuleComponent {
	return []*internal.ModuleComponent{}
}

func (m *SuggestionModule) Modals() []*internal.ModuleModal {
	return []*internal.ModuleModal{
		{"suggest:submit:", m.SuggestionSubmitModal},
	}
}

func (m *SuggestionModule) MenuCommands() map[string]internal.ModuleCommand {
	return map[string]internal.ModuleCommand{}
}

func New(db *database.Database, client *internal.MPBot, config *config.Config) *SuggestionModule {
	return &SuggestionModule{
		db:     db,
		client: client,
		config: config,
	}
}
