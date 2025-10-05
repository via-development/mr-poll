package suggestionModule

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/via-development/mr-poll/bot/internal/database"
	moduleUtil "github.com/via-development/mr-poll/bot/internal/util/module"
)

type SuggestionModule struct {
	moduleUtil.Module

	db     *database.GormDB
	client bot.Client
}

func (m *SuggestionModule) Name() string {
	return "suggestion"
}

func (m *SuggestionModule) Commands() map[string]moduleUtil.ModuleCommand {
	return map[string]moduleUtil.ModuleCommand{
		"suggest":    m.SuggestCommand,
		"suggestion": m.SuggestionCommand,
	}
}

func (m *SuggestionModule) Buttons() []*moduleUtil.ModuleComponent {
	return []*moduleUtil.ModuleComponent{
		{"suggestions:", m.SuggestionsVoteButton},
	}
}

func (m *SuggestionModule) SelectMenus() []*moduleUtil.ModuleComponent {
	return []*moduleUtil.ModuleComponent{}
}

func (m *SuggestionModule) Modals() []*moduleUtil.ModuleModal {
	return []*moduleUtil.ModuleModal{
		{"suggest:submit:", m.SuggestionSubmitModal},
	}
}

func (m *SuggestionModule) MenuCommands() map[string]moduleUtil.ModuleCommand {
	return map[string]moduleUtil.ModuleCommand{}
}

func New(db *database.GormDB, client bot.Client) *SuggestionModule {
	return &SuggestionModule{
		db:     db,
		client: client,
	}
}
