package util

import (
	"github.com/disgoorg/disgo/events"
)

type Module struct {
	Name        string
	Commands    map[string]ModuleCommand
	SelectMenus []*ModuleComponent
	Buttons     []*ModuleComponent
	Modals      []*ModuleModal
}

type ModuleCommand func(*events.ApplicationCommandInteractionCreate) error

type ModuleComponent struct {
	Prefix  string
	Execute func(*events.ComponentInteractionCreate) error
}

type ModuleModal struct {
	Prefix  string
	Execute func(create *events.ModalSubmitInteractionCreate) error
}
