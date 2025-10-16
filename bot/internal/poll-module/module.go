package poll_module

import (
	"github.com/via-development/mr-poll/bot/internal"
	"github.com/via-development/mr-poll/bot/internal/database"
	"go.uber.org/zap"
)

type PollModule struct {
	internal.Module

	db     *database.GormDB
	client *internal.MPBot
	log    *zap.Logger
}

func (m *PollModule) Name() string {
	return "poll"
}

func (m *PollModule) Commands() map[string]internal.ModuleCommand {
	return map[string]internal.ModuleCommand{
		"poll": m.PollCommand,
	}
}

func (m *PollModule) Buttons() []*internal.ModuleComponent {
	return []*internal.ModuleComponent{
		{"poll:option-", m.PollOptionButton},
		{"poll:menu", m.PollMenuButton},
	}
}

func (m *PollModule) SelectMenus() []*internal.ModuleComponent {
	return []*internal.ModuleComponent{
		{"poll:opts", m.PollOptionSelectMenu},
	}
}

func (m *PollModule) Modals() []*internal.ModuleModal {
	return []*internal.ModuleModal{
		{"poll:option-submit", m.PollOptionSubmitModal},
	}
}

func (m *PollModule) MenuCommands() map[string]internal.ModuleCommand {
	return map[string]internal.ModuleCommand{
		"End poll":     m.MenuPollEndCommand,
		"Refresh poll": m.MenuPollRefreshCommand,
	}
}

func New(db *database.GormDB, client *internal.MPBot, log *zap.Logger) *PollModule {
	return &PollModule{
		db:     db,
		client: client,
		log:    log,
	}
}
