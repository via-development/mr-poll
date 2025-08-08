package pollModule

import (
	pollButtons "github.com/via-development/mr-poll/bot/internal/poll-module/buttons"
	pollCommands "github.com/via-development/mr-poll/bot/internal/poll-module/commands"
	pollModals "github.com/via-development/mr-poll/bot/internal/poll-module/modals"
	pollSelectMenus "github.com/via-development/mr-poll/bot/internal/poll-module/select-menus"
	moduleUtil "github.com/via-development/mr-poll/bot/internal/util/module"
)

type PollModule struct {
	moduleUtil.Module
}

func (m *PollModule) Name() string {
	return "poll"
}

func (m *PollModule) Commands() map[string]moduleUtil.ModuleCommand {
	return map[string]moduleUtil.ModuleCommand{
		"poll": pollCommands.PollCommand,
	}
}

func (m *PollModule) Buttons() []*moduleUtil.ModuleComponent {
	return []*moduleUtil.ModuleComponent{
		{"poll:option-", pollButtons.PollOptionButton},
		{"poll:menu", pollButtons.PollMenuButton},
	}
}

func (m *PollModule) SelectMenus() []*moduleUtil.ModuleComponent {
	return []*moduleUtil.ModuleComponent{
		{"poll:opts-", pollSelectMenus.PollOptionSelectMenu},
	}
}

func (m *PollModule) Modals() []*moduleUtil.ModuleModal {
	return []*moduleUtil.ModuleModal{
		{"poll:option-submit", pollModals.PollOptionSubmitModal},
	}
}

//var Module = &util.Module{
//	Name: "poll",
//	Commands: map[string]util.ModuleCommand{
//		"poll": pollCommands.PollCommand,
//	},
//	Buttons: []*util.ModuleComponent{
//		{"poll:option-", pollButtons.PollOptionButton},
//		{"poll:menu", pollButtons.PollMenuButton},
//	},
//	SelectMenus: []*util.ModuleComponent{
//		{"poll:opts", pollSelectMenus.PollOptionSelectMenu},
//	},
//	Modals: []*util.ModuleModal{
//		{"poll:option-submit", pollModals.PollOptionSubmitModal},
//	},
//}

func New() *PollModule {
	return &PollModule{}
}
