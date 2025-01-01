package baseUtil

import (
	"github.com/disgoorg/disgo/events"
)

type Module struct {
	Commands    map[string]ModuleCommand
	SelectMenus []*ModuleComponent
	Buttons     []*ModuleComponent
}

type ModuleCommand func(*events.ApplicationCommandInteractionCreate) error

type ModuleComponent struct {
	Prefix  string
	Execute func(*events.ComponentInteractionCreate) error
}
