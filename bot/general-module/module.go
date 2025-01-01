package generalModule

import (
	baseUtil "mrpoll_bot/base-util"
	generalCommands "mrpoll_bot/general-module/commands"
	generalSelectMenus "mrpoll_bot/general-module/select-menus"
)

var Module = &baseUtil.Module{
	Commands: map[string]baseUtil.ModuleCommand{
		"help":    generalCommands.MrPollCommand,
		"mr-poll": generalCommands.MrPollCommand,
	},
	SelectMenus: []*baseUtil.ModuleComponent{
		{"help:", generalSelectMenus.MrPollSelectMenu},
	},
}
