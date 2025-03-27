package pollButtons

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"mrpoll_bot/database"
	pollUtil "mrpoll_bot/poll-module/util"
	"mrpoll_bot/util"
	"strconv"
)

func PollOptionButton(interaction *events.ComponentInteractionCreate) error {
	var pollData database.PollData
	if err := database.DB.Preload("Options").First(&pollData, interaction.Message.ID.String()).Error; err != nil {
		_ = interaction.CreateMessage(pollUtil.PollNotFoundMessage())
		panic(err)
		return err
	}

	var selectedOptions []int
	{
		s := interaction.Data.CustomID()[len("poll:option-"):]
		if !util.NumRegex.Match([]byte(s)) {
			// TODO: migrate poll
			return nil
		}
		n, _ := strconv.Atoi(s)
		selectedOptions = []int{n}
	}

	userId := interaction.User().ID.String()

	action, err := pollUtil.VotePoll(pollData, userId, selectedOptions)
	fmt.Printf("%v\n", pollData)
	if err != nil {
		return err
	}

	interaction.UpdateMessage(discord.MessageUpdate{
		Embeds: &[]discord.Embed{
			pollUtil.MakePollEmbed(pollData),
		},
	})

	interaction.Client().Rest().CreateFollowupMessage(interaction.Client().ID(), interaction.Token(), discord.MessageCreate{
		Flags:   discord.MessageFlagEphemeral,
		Content: fmt.Sprintf("Your vote was %s", action),
	})
	return nil
}
