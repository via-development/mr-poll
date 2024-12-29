package eventHandlers

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// ComponentTypeActionRow ComponentType = iota + 1
// ComponentTypeButton
// ComponentTypeStringSelectMenu
// ComponentTypeTextInput
// ComponentTypeUserSelectMenu
// ComponentTypeRoleSelectMenu
// ComponentTypeMentionableSelectMenu
// ComponentTypeChannelSelectMenu

// ComponentHandler Passes component interaction events to be handled by the appropriate functions
func ComponentHandler(e *events.ComponentInteractionCreate) {
	switch e.Data.Type() {
	case discord.ComponentTypeStringSelectMenu, discord.ComponentTypeUserSelectMenu, discord.ComponentTypeRoleSelectMenu, discord.ComponentTypeChannelSelectMenu:
		SelectMenuHandler(e)
		return
	case discord.ComponentTypeButton:
		return
	default: // ComponentTypeTextInput, ComponentTypeActionRow, ComponentTypeMentionableSelectMenu
		return
	}
}
