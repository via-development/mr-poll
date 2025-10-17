package poll_module

import (
	"github.com/via-development/mr-poll/bot/internal/config"
	"github.com/via-development/mr-poll/bot/internal/core"
	"github.com/via-development/mr-poll/bot/internal/database"
	"go.uber.org/zap"
)

type PollModule struct {
	core.Module

	db     *database.Database
	client *core.Client
	log    *zap.Logger
	config *config.Config
}

func (m *PollModule) Name() string {
	return "poll"
}

func (m *PollModule) Commands() map[string]core.ModuleCommand {
	return map[string]core.ModuleCommand{
		"poll": m.PollCommand,
	}
}

func (m *PollModule) Buttons() []*core.ModuleComponent {
	return []*core.ModuleComponent{
		{"poll:option-", m.PollOptionButton},
		{"poll:menu", m.PollMenuButton},
	}
}

func (m *PollModule) SelectMenus() []*core.ModuleComponent {
	return []*core.ModuleComponent{
		{"poll:opts", m.PollOptionSelectMenu},
	}
}

func (m *PollModule) Modals() []*core.ModuleModal {
	return []*core.ModuleModal{
		{"poll:option-submit", m.PollOptionSubmitModal},
	}
}

func (m *PollModule) MenuCommands() map[string]core.ModuleCommand {
	return map[string]core.ModuleCommand{
		"End poll":     m.MenuPollEndCommand,
		"Refresh poll": m.MenuPollRefreshCommand,
	}
}

func New(db *database.Database, client *core.Client, log *zap.Logger, config *config.Config) *PollModule {
	return &PollModule{
		db:     db,
		client: client,
		log:    log,
		config: config,
	}
}
