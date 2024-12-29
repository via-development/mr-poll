package generalModule

import (
	"github.com/disgoorg/disgo/bot"
	baseUtil "mrpoll_bot/base-util"
	generalCommands "mrpoll_bot/general-module/commands"
	generalSelectMenus "mrpoll_bot/general-module/select-menus"
)

type GeneralModule baseUtil.Module

var Module *GeneralModule

func InitGeneralModule(client *bot.Client) {
	Module = &GeneralModule{
		Client: client,
		Commands: map[string]baseUtil.ModuleCommand{
			"help":    generalCommands.MrPollCommand,
			"mr-poll": generalCommands.MrPollCommand,
		},
		SelectMenus: []*baseUtil.ModuleSelectMenu{
			{"help:", generalSelectMenus.MrPollSelectMenu},
		},
	}
}
