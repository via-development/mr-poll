package poll_module

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	"github.com/via-development/mr-poll/bot/internal/util"
)

func (m *PollModule) PollOptionButton(interaction *events.ComponentInteractionCreate) error {
	var poll schema.Poll
	if err := m.db.Preload("Options").First(&poll, interaction.Message.ID.String()).Error; err != nil {
		return interaction.CreateMessage(PollNotFoundMessage())
	}
	if poll.HasEnded {
		m.client.Rest.CreateFollowupMessage(interaction.Client().ID(), interaction.Token(), discord.MessageCreate{
			Flags:   discord.MessageFlagEphemeral,
			Content: "This poll has ended lol!",
		})
	}

	err := m.FetchPollUser(&poll)
	if err != nil {
		return err
	}

	var selectedOptions []int
	{
		s := interaction.Data.CustomID()[len("poll:option-"):]
		if s == "submit" {
			return interaction.Modal(PollOptionSubmitModel(interaction.Message.ID.String()))
		}
		if !util.NumRegex.Match([]byte(s)) {
			// Migrates Poll
			pollEmbeds := m.MakePollEmbeds(&poll)
			components := m.MakePollComponents(&poll)
			err = interaction.UpdateMessage(discord.MessageUpdate{
				Embeds:     &pollEmbeds,
				Components: &components,
			})
			m.client.Rest.CreateFollowupMessage(interaction.Client().ID(), interaction.Token(), discord.MessageCreate{
				Flags:   discord.MessageFlagEphemeral,
				Content: "Hey, I have just updated this poll to work with the new system, please vote again!",
			})
			return err
		}
		n, _ := strconv.Atoi(s)
		selectedOptions = []int{n}
	}

	userId := interaction.User().ID.String()

	action, err := m.VotePoll(&poll, userId, selectedOptions)
	if err != nil {
		return err
	}

	pollEmbeds := m.MakePollEmbeds(&poll)
	messageUpdate := discord.MessageUpdate{
		Embeds: &pollEmbeds,
	}
	if poll.CanSubmit {
		pollComponents := m.MakePollComponents(&poll)
		messageUpdate.Components = &pollComponents
	}
	err = interaction.UpdateMessage(messageUpdate)

	interaction.Client().Rest.CreateFollowupMessage(interaction.Client().ID(), interaction.Token(), discord.MessageCreate{
		Flags:  discord.MessageFlagEphemeral,
		Embeds: []discord.Embed{util.MakeSuccessEmbed("Your vote was " + action)},
	})
	return err
}

func (m *PollModule) PollMenuButton(interaction *events.ComponentInteractionCreate) error {
	if !strings.HasPrefix(interaction.Data.CustomID(), "poll:menu-") {
		return m.handlePollMenuMainPage(interaction)
	}
	parts := strings.Split(interaction.Data.CustomID()[len("poll:menu-"):], ":")
	op, pollId := parts[0], parts[1]

	var poll schema.Poll
	res := m.db.Preload("Options").First(&poll, pollId)
	if res.Error != nil {
		return interaction.CreateMessage(PollNotFoundMessage())
	}

	switch op {
	case "refresh":
		if poll.GuildId == nil {
			return nil
		}

		err := m.FetchPollUser(&poll)
		if err != nil {
			return err
		}

		e := m.MakePollEmbeds(&poll)
		c := m.MakePollComponents(&poll)

		_, err = interaction.Client().Rest.UpdateMessage(poll.ChannelIdSnowflake(), poll.MessageIdSnowflake(), discord.MessageUpdate{
			Embeds:     &e,
			Components: &c,
		})
		if err != nil {
			return err
		}

		return interaction.CreateMessage(discord.MessageCreate{
			Flags:  discord.MessageFlagEphemeral,
			Embeds: []discord.Embed{util.MakeSuccessEmbed("The poll has been refreshed!")},
		})
	case "end":
		// If it's not their own poll, check for permissions
		if poll.UserId != interaction.User().ID.String() {
			if poll.GuildId == nil || interaction.GuildID() == nil {
				return interaction.CreateMessage(NotYourPollMessage())
			}

			channel, _ := interaction.Client().Caches.GuildTextChannel(interaction.Channel().ID())
			perms := interaction.Client().Caches.MemberPermissionsInChannel(channel, interaction.Member().Member)

			if !perms.Has(discord.PermissionManageMessages) {
				return interaction.CreateMessage(NotYourPollMessage())
			}
		}

		err := m.FetchPollUser(&poll)
		if err != nil {
			return err
		}

		enderId := interaction.User().ID.String()

		err = m.EndPoll(&poll, &enderId)
		if err != nil {
			return err
		}

		if poll.GuildId == nil {

		} else {
			err = m.UpdatePollMessage(&poll)
		}

		if err != nil {
			return err
		}

		return interaction.CreateMessage(discord.MessageCreate{
			Flags: discord.MessageFlagEphemeral,
			Embeds: []discord.Embed{
				util.MakeSuccessEmbed(fmt.Sprintf("Ended [poll](%s)!", poll.MessageUrl())),
			},
		})
	}

	return nil
}

func (m *PollModule) handlePollMenuMainPage(interaction *events.ComponentInteractionCreate) error {
	var pollData schema.Poll
	res := m.db.Preload("Options").First(&pollData, interaction.Message.ID.String())
	if res.Error != nil {
		return interaction.CreateMessage(PollNotFoundMessage())
	}

	buttons := discord.ActionRowComponent{
		Components: []discord.InteractiveComponent{
			discord.ButtonComponent{
				Style:    discord.ButtonStyleSecondary,
				Emoji:    &discord.ComponentEmoji{Name: "🔄", ID: 0},
				CustomID: "poll:menu-refresh:" + pollData.MessageId,
			},
			discord.ButtonComponent{
				Style:    discord.ButtonStyleSecondary,
				Emoji:    &discord.ComponentEmoji{Name: "🛑", ID: 0},
				CustomID: "poll:menu-end:" + pollData.MessageId,
			},
		},
	}
	return interaction.CreateMessage(discord.MessageCreate{
		Flags:      discord.MessageFlagEphemeral,
		Content:    fmt.Sprintf("ℹ️ You can do the same actions by right clicking the poll!\n> %s", pollData.Question),
		Components: []discord.LayoutComponent{buttons},
	})
}
