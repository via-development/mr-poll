package baseUtil

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
)

type Module struct {
	Client      *bot.Client
	Commands    map[string]ModuleCommand
	SelectMenus []*ModuleSelectMenu
}

type ModuleCommand func(*events.ApplicationCommandInteractionCreate) error

type ModuleSelectMenu struct {
	Prefix  string
	Execute func(*events.ComponentInteractionCreate) error
}
