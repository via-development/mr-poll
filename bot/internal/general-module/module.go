package generalModule

import (
	"github.com/via-development/mr-poll/bot/internal"
	"github.com/via-development/mr-poll/bot/internal/api"
	"github.com/via-development/mr-poll/bot/internal/config"
)

type GeneralModule struct {
	internal.Module

	api    *api.Api
	config *config.Config
}

func (m *GeneralModule) Name() string {
	return "general"
}

func (m *GeneralModule) Commands() map[string]internal.ModuleCommand {
	return map[string]internal.ModuleCommand{
		"help":    m.MrPollCommand,
		"mr-poll": m.MrPollCommand,
		"mytime":  m.MyTimeCommand,
	}
}

func (m *GeneralModule) Buttons() []*internal.ModuleComponent {
	return []*internal.ModuleComponent{
		{"help:", m.MrPollButton},
	}
}

func (m *GeneralModule) SelectMenus() []*internal.ModuleComponent {
	return []*internal.ModuleComponent{}
}

func (m *GeneralModule) Modals() []*internal.ModuleModal {
	return []*internal.ModuleModal{}
}

func (m *GeneralModule) MenuCommands() map[string]internal.ModuleCommand {
	return map[string]internal.ModuleCommand{}
}

func New(api *api.Api, config *config.Config) *GeneralModule {
	return &GeneralModule{
		api:    api,
		config: config,
	}
}
