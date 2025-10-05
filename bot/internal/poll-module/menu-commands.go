package poll_module

import (
	"errors"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	"gorm.io/gorm"
)

func (m *PollModule) MenuPollEndCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	var poll schema.Poll
	err := m.db.Preload("Options").Find(&poll, schema.Poll{
		MessageId: interaction.MessageCommandInteractionData().TargetMessage().ID.String(),
	}).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return interaction.CreateMessage(PollNotFoundMessage())
	} else if err != nil {
		return err
	}

	if poll.UserId != interaction.User().ID.String() {
		if interaction.GuildID() == nil || poll.GuildId == nil {
			return interaction.CreateMessage(NotYourPollMessage())
		}

		channel, _ := interaction.Client().Caches().GuildTextChannel(interaction.Channel().ID())
		perms := interaction.Client().Caches().MemberPermissionsInChannel(channel, interaction.Member().Member)

		if !perms.Has(discord.PermissionManageMessages) {
			return interaction.CreateMessage(NotYourPollMessage())
		}
	}

	err = m.FetchPollUser(&poll)
	if err != nil {
		return err
	}

	enderId := interaction.User().ID.String()

	err = m.EndPoll(&poll, &enderId)
	if err != nil {
		return err
	}

	content := MakePollText(&poll)
	pollEmbeds := m.MakePollEmbeds(&poll)
	pollComponents := m.MakePollComponents(&poll)

	messageUpdate := discord.MessageUpdate{
		Content:    &content,
		Embeds:     &pollEmbeds,
		Components: &pollComponents,
	}

	_, err = interaction.Client().Rest().UpdateMessage(interaction.Channel().ID(), interaction.MessageCommandInteractionData().TargetMessage().ID, messageUpdate)
	return err
}

func (m *PollModule) MenuPollRefreshCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	return nil
}
