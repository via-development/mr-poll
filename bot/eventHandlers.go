package main

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"mrpoll_bot/database"
	generalModule "mrpoll_bot/general-module"
	pollModule "mrpoll_bot/poll-module"
	suggestionModule "mrpoll_bot/suggestion-module"
	"mrpoll_bot/util"
	"slices"
	"strings"
)

var modules = []*util.Module{
	generalModule.Module,
	pollModule.Module,
	suggestionModule.Module,
}

func CommandHandler(e *events.ApplicationCommandInteractionCreate) {
	commandName := e.SlashCommandInteractionData().CommandName()

	for _, module := range modules {
		command, ok := module.Commands[commandName]
		if ok {
			if slices.Index(database.BotSettingsC.DisabledModules, module.Name) != -1 {
				e.CreateMessage(util.DisabledModuleMessage())
				return
			}
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
				if slices.Index(database.BotSettingsC.DisabledModules, module.Name) != -1 {
					e.CreateMessage(util.DisabledModuleMessage())
					return
				}
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
				if slices.Index(database.BotSettingsC.DisabledModules, module.Name) != -1 {
					e.CreateMessage(util.DisabledModuleMessage())
					return
				}
				_ = selectMenu.Execute(e)
				return
			}
		}
	}
}

func ModalHandler(e *events.ModalSubmitInteractionCreate) {
	customId := e.Data.CustomID

	for _, module := range modules {
		for _, modal := range module.Modals {
			if strings.HasPrefix(customId, modal.Prefix) {
				if slices.Index(database.BotSettingsC.DisabledModules, module.Name) != -1 {
					e.CreateMessage(util.DisabledModuleMessage())
					return
				}
				_ = modal.Execute(e)
				return
			}
		}
	}
}

func MessageHandler(e *events.MessageCreate) {
	if !strings.HasPrefix(e.Message.Content, "<@"+e.Client().ID().String()+"> ") {
		return
	}
	args := strings.Split(e.Message.Content[len("<@"+e.Client().ID().String()+"> "):], " ")
	command := args[0]
	args = args[1:]

	switch command {
	case "disable":
		{
			if len(args) > 0 {
				if args[0] == "all" {
					l := len(database.BotSettingsC.DisabledModules)
					database.BotSettingsC.DisabledModules = []string{}
					if l != len(modules) {
						for _, module := range modules {
							database.BotSettingsC.DisabledModules = append(database.BotSettingsC.DisabledModules, module.Name)
						}
					}
				} else {
					for _, module := range args {
						if i := slices.Index(database.BotSettingsC.DisabledModules, module); i != -1 {
							database.BotSettingsC.DisabledModules = append(database.BotSettingsC.DisabledModules[:i], database.BotSettingsC.DisabledModules[i+1:]...)
						} else {
							database.BotSettingsC.DisabledModules = append(database.BotSettingsC.DisabledModules, module)
						}
					}
				}

				database.DB.Save(database.BotSettingsC)
			}

			e.Client().Rest().CreateMessage(e.ChannelID, discord.MessageCreate{
				Content:          "Disabled: " + strings.Join(database.BotSettingsC.DisabledModules, ", "),
				MessageReference: e.Message.MessageReference,
			})
		}

	}
}
