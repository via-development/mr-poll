package pollButtons

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/via-development/mr-poll/bot/internal/database"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	pollUtil "github.com/via-development/mr-poll/bot/internal/poll-module/util"
	"github.com/via-development/mr-poll/bot/internal/util"
	"strings"
)

func PollMenuButton(interaction *events.ComponentInteractionCreate, db *database.GormDB) error {
	if !strings.HasPrefix(interaction.Data.CustomID(), "poll:menu-") {
		return handlePollMenuMainPage(interaction, db)
	}
	parts := strings.Split(interaction.Data.CustomID()[len("poll:menu-"):], ":")
	op, pollId := parts[0], parts[1]

	var pollData schema.Poll
	res := db.Preload("Options").First(&pollData, pollId)
	if res.Error != nil {
		return interaction.CreateMessage(pollUtil.PollNotFoundMessage())
	}

	switch op {
	case "refresh":
		if pollData.GuildId == nil {
			return nil
		}

		err := pollUtil.FetchPollUser(interaction.Client(), db, &pollData)
		if err != nil {
			return err
		}

		e := pollUtil.MakePollEmbeds(&pollData)
		c := pollUtil.MakePollComponents(interaction.Client(), db, &pollData)

		_, err = interaction.Client().Rest().UpdateMessage(pollData.ChannelIdSnowflake(), pollData.MessageIdSnowflake(), discord.MessageUpdate{
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
	}
	return nil
}

func handlePollMenuMainPage(interaction *events.ComponentInteractionCreate, db *database.GormDB) error {
	var pollData schema.Poll
	res := db.Preload("Options").First(&pollData, interaction.Message.ID.String())
	if res.Error != nil {
		return interaction.CreateMessage(pollUtil.PollNotFoundMessage())
	}

	buttons := discord.ActionRowComponent{
		discord.ButtonComponent{
			Style:    discord.ButtonStyleSecondary,
			Emoji:    &discord.ComponentEmoji{Name: "🔄", ID: 0},
			CustomID: "poll:menu-refresh:" + pollData.MessageId,
		},
	}
	return interaction.CreateMessage(discord.MessageCreate{
		Flags:      discord.MessageFlagEphemeral,
		Content:    fmt.Sprintf("ℹ️ You can do the same actions by right clicking the poll!\n> %s", pollData.Question),
		Components: []discord.ContainerComponent{buttons},
	})
}
