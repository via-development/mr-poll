package poll_module

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/via-development/mr-poll/bot/internal/database"
	moduleUtil "github.com/via-development/mr-poll/bot/internal/util/module"
	"go.uber.org/zap"
)

type PollModule struct {
	moduleUtil.Module

	db     *database.GormDB
	client bot.Client
	log    *zap.Logger
}

func (m *PollModule) Name() string {
	return "poll"
}

func (m *PollModule) Commands() map[string]moduleUtil.ModuleCommand {
	return map[string]moduleUtil.ModuleCommand{
		"poll": m.PollCommand,
	}
}

func (m *PollModule) Buttons() []*moduleUtil.ModuleComponent {
	return []*moduleUtil.ModuleComponent{
		{"poll:option-", m.PollOptionButton},
		{"poll:menu", m.PollMenuButton},
	}
}

func (m *PollModule) SelectMenus() []*moduleUtil.ModuleComponent {
	return []*moduleUtil.ModuleComponent{
		{"poll:opts", m.PollOptionSelectMenu},
	}
}

func (m *PollModule) Modals() []*moduleUtil.ModuleModal {
	return []*moduleUtil.ModuleModal{
		{"poll:option-submit", m.PollOptionSubmitModal},
	}
}

func (m *PollModule) MenuCommands() map[string]moduleUtil.ModuleCommand {
	return map[string]moduleUtil.ModuleCommand{
		"End poll":     m.MenuPollEndCommand,
		"Refresh poll": m.MenuPollRefreshCommand,
	}
}

func New(db *database.GormDB, client bot.Client, log *zap.Logger) *PollModule {
	return &PollModule{
		db:     db,
		client: client,
		log:    log,
	}
}
