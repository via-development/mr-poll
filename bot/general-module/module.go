package generalModule

import (
	generalCommands "mrpoll_bot/general-module/commands"
	generalSelectMenus "mrpoll_bot/general-module/select-menus"
	"mrpoll_bot/util"
)

var Module = &util.Module{
	Name: "general",
	Commands: map[string]util.ModuleCommand{
		"help":    generalCommands.MrPollCommand,
		"mr-poll": generalCommands.MrPollCommand,
	},
	SelectMenus: []*util.ModuleComponent{
		{"help:", generalSelectMenus.MrPollSelectMenu},
	},
}
