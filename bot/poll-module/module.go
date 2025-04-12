package pollModule

import (
	"mrpoll_bot/poll-module/buttons"
	pollCommands "mrpoll_bot/poll-module/commands"
	pollModals "mrpoll_bot/poll-module/modals"
	pollSelectMenus "mrpoll_bot/poll-module/select-menus"
	"mrpoll_bot/util"
)

var Module = &util.Module{
	Name: "poll",
	Commands: map[string]util.ModuleCommand{
		"poll": pollCommands.PollCommand,
	},
	Buttons: []*util.ModuleComponent{
		{"poll:option-", pollButtons.PollOptionButton},
		{"poll:menu", pollButtons.PollMenuButton},
	},
	SelectMenus: []*util.ModuleComponent{
		{"poll:opts", pollSelectMenus.PollOptionSelectMenu},
	},
	Modals: []*util.ModuleModal{
		{"poll:option-submit", pollModals.PollOptionSubmitModal},
	},
}
