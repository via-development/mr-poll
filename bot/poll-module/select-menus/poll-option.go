package pollSelectMenus

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"mrpoll_bot/database"
	pollUtil "mrpoll_bot/poll-module/util"
	"strconv"
)

func PollOptionSelectMenu(interaction *events.ComponentInteractionCreate) error {
	var pollData database.PollData
	if err := database.DB.Preload("Options").First(&pollData, interaction.Message.ID.String()).Error; err != nil {
		_ = interaction.CreateMessage(pollUtil.PollNotFoundMessage())
		return err
	}
	pollData.FetchUser(interaction.Client())
	user := interaction.User()
	database.DB.Save(&database.UserData{
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

	action, err := pollUtil.VotePoll(pollData, userId, selectedOptions)
	if err != nil {
		return err
	}

	pollEmbeds := pollUtil.MakePollEmbeds(pollData)
	err = interaction.UpdateMessage(discord.MessageUpdate{
		Embeds: &pollEmbeds,
	})

	_, _ = interaction.Client().Rest().CreateFollowupMessage(interaction.Client().ID(), interaction.Token(), discord.MessageCreate{
		Flags:   discord.MessageFlagEphemeral,
		Content: fmt.Sprintf("Your vote was %s", action),
	})
	return err
}
