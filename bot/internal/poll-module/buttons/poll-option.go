package pollButtons

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/via-development/mr-poll/bot/internal/database"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	pollUtil "github.com/via-development/mr-poll/bot/internal/poll-module/util"
	"github.com/via-development/mr-poll/bot/internal/util"
	"strconv"
)

func PollOptionButton(interaction *events.ComponentInteractionCreate, db *database.GormDB) error {
	var pollData schema.Poll
	if err := db.Preload("Options").First(&pollData, interaction.Message.ID.String()).Error; err != nil {
		return interaction.CreateMessage(pollUtil.PollNotFoundMessage())
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

	var selectedOptions []int
	{
		s := interaction.Data.CustomID()[len("poll:option-"):]
		if !util.NumRegex.Match([]byte(s)) {
			// Migrates Poll
			pollEmbeds := pollUtil.MakePollEmbeds(&pollData)
			components := pollUtil.MakePollComponents(interaction.Client(), db, &pollData)
			err = interaction.UpdateMessage(discord.MessageUpdate{
				Embeds:     &pollEmbeds,
				Components: &components,
			})
			interaction.Client().Rest().CreateFollowupMessage(interaction.Client().ID(), interaction.Token(), discord.MessageCreate{
				Flags:   discord.MessageFlagEphemeral,
				Content: "Hey, I have just updated this poll to work with the new system, please vote again!",
			})
			return err
		}
		n, _ := strconv.Atoi(s)
		selectedOptions = []int{n}
	}

	userId := interaction.User().ID.String()

	action, err := pollUtil.VotePoll(db, &pollData, userId, selectedOptions)
	if err != nil {
		return err
	}

	pollEmbeds := pollUtil.MakePollEmbeds(&pollData)
	err = interaction.UpdateMessage(discord.MessageUpdate{
		Embeds: &pollEmbeds,
	})

	interaction.Client().Rest().CreateFollowupMessage(interaction.Client().ID(), interaction.Token(), discord.MessageCreate{
		Flags:  discord.MessageFlagEphemeral,
		Embeds: []discord.Embed{util.MakeSuccessEmbed("Your vote was " + action)},
	})
	return err
}
