package generalModule

import (
	"github.com/via-development/mr-poll/bot/internal/config"
	"github.com/via-development/mr-poll/bot/internal/core"
	"github.com/via-development/mr-poll/bot/internal/redis"
)

type GeneralModule struct {
	core.Module

	redis  *redis.Client
	config *config.Config
}

func (m *GeneralModule) Name() string {
	return "general"
}

func (m *GeneralModule) Commands() map[string]core.ModuleCommand {
	return map[string]core.ModuleCommand{
		"help":    m.MrPollCommand,
		"mr-poll": m.MrPollCommand,
		"mytime":  m.MyTimeCommand,
	}
}

func (m *GeneralModule) Buttons() []*core.ModuleComponent {
	return []*core.ModuleComponent{
		{"help:", m.MrPollButton},
	}
}

func (m *GeneralModule) SelectMenus() []*core.ModuleComponent {
	return []*core.ModuleComponent{}
}

func (m *GeneralModule) Modals() []*core.ModuleModal {
	return []*core.ModuleModal{}
}

func (m *GeneralModule) MenuCommands() map[string]core.ModuleCommand {
	return map[string]core.ModuleCommand{}
}

func New(redis *redis.Client, config *config.Config) *GeneralModule {
	return &GeneralModule{
		redis:  redis,
		config: config,
	}
}
