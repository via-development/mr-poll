package generalModule

import (
	moduleUtil "github.com/via-development/mr-poll/bot/internal/util/module"
)

type GeneralModule struct {
	moduleUtil.Module
}

func (m *GeneralModule) Name() string {
	return "general"
}

func (m *GeneralModule) Commands() map[string]moduleUtil.ModuleCommand {
	return map[string]moduleUtil.ModuleCommand{
		"help":    m.MrPollCommand,
		"mr-poll": m.MrPollCommand,
		"mytime":  m.MyTimeCommand,
	}
}

func (m *GeneralModule) Buttons() []*moduleUtil.ModuleComponent {
	return []*moduleUtil.ModuleComponent{
		{"help:", m.MrPollButton},
	}
}

func (m *GeneralModule) SelectMenus() []*moduleUtil.ModuleComponent {
	return []*moduleUtil.ModuleComponent{}
}

func (m *GeneralModule) Modals() []*moduleUtil.ModuleModal {
	return []*moduleUtil.ModuleModal{}
}

func (m *GeneralModule) MenuCommands() map[string]moduleUtil.ModuleCommand {
	return map[string]moduleUtil.ModuleCommand{}
}

func New() *GeneralModule {
	return &GeneralModule{}
}
