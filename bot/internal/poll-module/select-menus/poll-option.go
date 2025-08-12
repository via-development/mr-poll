package pollSelectMenus

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/via-development/mr-poll/bot/internal/database"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	pollUtil "github.com/via-development/mr-poll/bot/internal/poll-module/util"
	"strconv"
)

func PollOptionSelectMenu(interaction *events.ComponentInteractionCreate, db *database.GormDB) error {
	var pollData schema.Poll
	if err := db.Preload("Options").First(&pollData, interaction.Message.ID.String()).Error; err != nil {
		_ = interaction.CreateMessage(pollUtil.PollNotFoundMessage())
		return err
	}
	if pollData.HasEnded {
		interaction.Client().Rest().CreateFollowupMessage(interaction.Client().ID(), interaction.Token(), discord.MessageCreate{
			Flags:   discord.MessageFlagEphemeral,
			Content: "This poll has ended lol!",
		})
	}

	err := pollUtil.FetchPollUser(interaction.Client(), db, &pollData)
	if err != nil {
		return err
	}

	user := interaction.User()
	db.Save(&schema.User{
		UserId:      user.ID.String(),
		Username:    user.Username,
		DisplayName: *user.GlobalName,
	})

	userId := interaction.User().ID.String()

	var selectedOptions []int
	{
		sv := interaction.StringSelectMenuInteractionData().Values

		for i := range sv {
			s := sv[i][len("option-"):]
			if s == "submit" {
				return interaction.Modal(pollUtil.PollOptionSubmitModel(interaction.Message.ID.String()))
			}
			n, _ := strconv.Atoi(s)
			selectedOptions = append(selectedOptions, n)
		}
	}

	action, err := pollUtil.VotePoll(db, &pollData, userId, selectedOptions)
	if err != nil {
		return err
	}

	pollEmbeds := pollUtil.MakePollEmbeds(&pollData)
	err = interaction.UpdateMessage(discord.MessageUpdate{
		Embeds: &pollEmbeds,
	})

	_, _ = interaction.Client().Rest().CreateFollowupMessage(interaction.Client().ID(), interaction.Token(), discord.MessageCreate{
		Flags:   discord.MessageFlagEphemeral,
		Content: fmt.Sprintf("Your vote was %s", action),
	})
	return err
}
