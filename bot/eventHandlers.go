package main

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	baseUtil "mrpoll_bot/base-util"
	generalModule "mrpoll_bot/general-module"
	pollModule "mrpoll_bot/poll-module"
	suggestionModule "mrpoll_bot/suggestion-module"
	"strings"
)

var modules = []*baseUtil.Module{
	generalModule.Module,
	pollModule.Module,
	suggestionModule.Module,
}

func CommandHandler(e *events.ApplicationCommandInteractionCreate) {
	commandName := e.SlashCommandInteractionData().CommandName()

	for _, module := range modules {
		command, ok := module.Commands[commandName]
		if ok {
			err := command(e)
			fmt.Println("Err: ", err)

			return
		}
	}

	_ = e.CreateMessage(discord.MessageCreate{
		Content: "I couldn't find that command!",
		Flags:   discord.MessageFlagEphemeral,
	})
}

// ComponentHandler Passes component interaction events to be handled by the appropriate functions
func ComponentHandler(e *events.ComponentInteractionCreate) {
	switch e.Data.Type() {
	case discord.ComponentTypeStringSelectMenu, discord.ComponentTypeUserSelectMenu, discord.ComponentTypeRoleSelectMenu, discord.ComponentTypeChannelSelectMenu:
		SelectMenuHandler(e)
		return
	case discord.ComponentTypeButton:
		ButtonHandler(e)
		return
	default: // ComponentTypeTextInput, ComponentTypeActionRow, ComponentTypeMentionableSelectMenu
		return
	}
}

func ButtonHandler(e *events.ComponentInteractionCreate) {
	customId := e.Data.CustomID()

	for _, module := range modules {
		for _, button := range module.Buttons {
			if strings.HasPrefix(customId, button.Prefix) {
				_ = button.Execute(e)
				return
			}
		}
	}
}

// SelectMenuHandler is not directly emitted by Disgo but by ComponentHandler
func SelectMenuHandler(e *events.ComponentInteractionCreate) {
	customId := e.Data.CustomID()

	for _, module := range modules {
		for _, selectMenu := range module.SelectMenus {
			if strings.HasPrefix(customId, selectMenu.Prefix) {
				_ = selectMenu.Execute(e)
				return
			}
		}
	}
}
