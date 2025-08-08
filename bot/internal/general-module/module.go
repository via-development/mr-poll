package generalModule

import (
	generalCommands "github.com/via-development/mr-poll/bot/general-module/commands"
	generalSelectMenus "github.com/via-development/mr-poll/bot/general-module/select-menus"
	"github.com/via-development/mr-poll/bot/util"
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
