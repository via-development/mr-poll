package core

import "github.com/disgoorg/disgo/events"

type Module interface {
	Name() string
	Commands() map[string]ModuleCommand
	SelectMenus() []*ModuleComponent
	Buttons() []*ModuleComponent
	Modals() []*ModuleModal
	MenuCommands() map[string]ModuleCommand
}

type ModuleCommand func(*events.ApplicationCommandInteractionCreate) error

type ModuleComponent struct {
	Prefix  string
	Execute func(*events.ComponentInteractionCreate) error
}

type ModuleModal struct {
	Prefix  string
	Execute func(*events.ModalSubmitInteractionCreate) error
}

type ModuleApp func(*events.ApplicationCommandInteractionCreate) error
