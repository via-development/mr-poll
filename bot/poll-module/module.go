package pollModule

import (
	"mrpoll_bot/poll-module/buttons"
	pollCommands "mrpoll_bot/poll-module/commands"
	pollSelectMenus "mrpoll_bot/poll-module/select-menus"
	"mrpoll_bot/util"
)

var Module = &util.Module{
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
}
