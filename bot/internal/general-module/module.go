package generalModule

import (
	generalCommands "github.com/via-development/mr-poll/bot/internal/general-module/commands"
	generalSelectMenus "github.com/via-development/mr-poll/bot/internal/general-module/select-menus"
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
		"help":    generalCommands.MrPollCommand,
		"mr-poll": generalCommands.MrPollCommand,
	}
}

func (m *GeneralModule) Buttons() []*moduleUtil.ModuleComponent {
	return []*moduleUtil.ModuleComponent{}
}

func (m *GeneralModule) SelectMenus() []*moduleUtil.ModuleComponent {
	return []*moduleUtil.ModuleComponent{
		{"help:", generalSelectMenus.MrPollSelectMenu},
	}
}

func (m *GeneralModule) Modals() []*moduleUtil.ModuleModal {
	return []*moduleUtil.ModuleModal{}
}

func New() *GeneralModule {
	return &GeneralModule{}
}
