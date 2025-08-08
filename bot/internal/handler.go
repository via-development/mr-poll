package internal

import (
	"context"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/via-development/mr-poll/bot/internal/util"
	"go.uber.org/zap"
	"strings"
)

func (b *MPBot) isModuleDisabled(module string) bool {
	for _, m := range b.db.BotSettings().DisabledModules {
		if m == module {
			return true
		}
	}
	return false
}

func (b *MPBot) HandleReady(e *events.Ready) {
	b.log.Info(fmt.Sprintf("shard %d is online", e.ShardID()))
	err := e.Client().SetPresenceForShard(
		context.Background(), e.ShardID(),
		gateway.WithCustomActivity(
			"online",
			gateway.WithActivityState(fmt.Sprintf("/mr-poll | %s Shard (%d)", util.ShardNames[e.ShardID()], e.ShardID())),
		),
	)

	if err != nil {
		b.log.Error("Could not set presence for shard", zap.Int("shardId", e.ShardID()), zap.Error(err))
	}
}

func (b *MPBot) HandleCommandInteraction(e *events.ApplicationCommandInteractionCreate) {
	commandName := e.SlashCommandInteractionData().CommandName()

	for _, module := range b.modules {
		command, ok := module.Commands()[commandName]
		if ok {
			if b.isModuleDisabled(module.Name()) {
				_ = e.CreateMessage(util.DisabledModuleMessage())
				return
			}

			err := command(e, b.db)
			if err != nil {
				b.log.Error("Failed to execute command!", zap.String("name", e.Data.CommandName()), zap.Error(err))
				return
			}

			return
		}
	}

	_ = e.CreateMessage(discord.MessageCreate{
		Content: "I couldn't find that command!",
		Flags:   discord.MessageFlagEphemeral,
	})
}

func (b *MPBot) HandleComponentInteraction(e *events.ComponentInteractionCreate) {
	switch e.Data.Type() {
	case discord.ComponentTypeStringSelectMenu, discord.ComponentTypeUserSelectMenu, discord.ComponentTypeRoleSelectMenu, discord.ComponentTypeChannelSelectMenu:
		b.HandleSelectMenuInteraction(e)
		return
	case discord.ComponentTypeButton:
		b.HandleButtonInteraction(e)
		return
	default: // ComponentTypeTextInput, ComponentTypeActionRow, ComponentTypeMentionableSelectMenu
		return
	}
}

func (b *MPBot) HandleButtonInteraction(e *events.ComponentInteractionCreate) {
	customId := e.Data.CustomID()

	for _, module := range b.modules {
		for _, button := range module.Buttons() {
			if strings.HasPrefix(customId, button.Prefix) {
				if b.isModuleDisabled(module.Name()) {
					_ = e.CreateMessage(util.DisabledModuleMessage())
					return
				}

				err := button.Execute(e, b.db)
				if err != nil {
					b.log.Error("Failed to execute button!", zap.String("customId", e.Data.CustomID()), zap.Error(err))
					return
				}

				return
			}
		}
	}
}

func (b *MPBot) HandleSelectMenuInteraction(e *events.ComponentInteractionCreate) {
	customId := e.Data.CustomID()

	for _, module := range b.modules {
		for _, selectMenu := range module.SelectMenus() {
			if strings.HasPrefix(customId, selectMenu.Prefix) {
				if b.isModuleDisabled(module.Name()) {
					_ = e.CreateMessage(util.DisabledModuleMessage())
					return
				}

				err := selectMenu.Execute(e, b.db)
				if err != nil {
					b.log.Error("Failed to execute select menu!", zap.String("customId", e.Data.CustomID()), zap.Error(err))
					return
				}

				return
			}
		}
	}
}
func (b *MPBot) HandleModalSubmitInteraction(e *events.ModalSubmitInteractionCreate) {

}

func (b *MPBot) HandleMessage(e *events.ModalSubmitInteractionCreate) {

}
