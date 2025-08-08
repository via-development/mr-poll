package moduleUtil

import (
	"github.com/disgoorg/disgo/events"
	"github.com/via-development/mr-poll/bot/internal/database"
)

type Module interface {
	Name() string
	Commands() map[string]ModuleCommand
	SelectMenus() []*ModuleComponent
	Buttons() []*ModuleComponent
	Modals() []*ModuleModal
}

type ModuleCommand func(*events.ApplicationCommandInteractionCreate, *database.GormDB) error

type ModuleComponent struct {
	Prefix  string
	Execute func(*events.ComponentInteractionCreate, *database.GormDB) error
}

type ModuleModal struct {
	Prefix  string
	Execute func(*events.ModalSubmitInteractionCreate, *database.GormDB) error
}
