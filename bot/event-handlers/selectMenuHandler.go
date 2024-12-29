package eventHandlers

import (
	"github.com/disgoorg/disgo/events"
	baseUtil "mrpoll_bot/base-util"
	generalModule "mrpoll_bot/general-module"
	pollModule "mrpoll_bot/poll-module"
	suggestionModule "mrpoll_bot/suggestion-module"
	"strings"
)

// SelectMenuHandler is not directly emitted by Disgo but by ComponentHandler
func SelectMenuHandler(e *events.ComponentInteractionCreate) {
	customId := e.Data.CustomID()
	modules := []*baseUtil.Module{
		(*baseUtil.Module)(generalModule.Module),
		(*baseUtil.Module)(pollModule.Module),
		(*baseUtil.Module)(suggestionModule.Module),
	}

	for _, module := range modules {
		for _, selectMenu := range module.SelectMenus {
			if strings.HasPrefix(customId, selectMenu.Prefix) {
				_ = selectMenu.Execute(e)
				return
			}
		}
	}

}
