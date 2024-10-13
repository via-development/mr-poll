package structs

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
)

type BotModule struct {
	Name     string
	Client   *bot.Client
	Commands map[string]func(module *BotModule, interaction *events.ApplicationCommandInteractionCreate)
}
