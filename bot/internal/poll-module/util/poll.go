package pollUtil

import (
	"errors"
	"fmt"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/via-development/mr-poll/bot/internal/database"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	"go.uber.org/zap"
	"slices"
	"time"
)

var requiredPermissions = discord.PermissionSendMessages // | discord.PermissionViewChannel

// CreatePoll is used by both the website and poll command for creating polls, it checks permissions, creates poll message and saves to database.
func CreatePoll(client bot.Client, db *database.GormDB, pollData *schema.PollData) (*discord.Message, error) {
	channelId := pollData.ChannelIdSnowflake()

	// Channel permission check
	channel, found := client.Caches().Channel(channelId)
	if !found {
		return nil, errors.New("channel is not in cache")
	}
	member, found := client.Caches().Member(*pollData.GuildIdSnowflake(), client.ApplicationID())
	if !found {
		m, err := client.Rest().GetMember(*pollData.GuildIdSnowflake(), client.ApplicationID())
		if err != nil {
			return nil, err
		}
		member = *m
	}
	p := client.Caches().MemberPermissionsInChannel(channel, member)
	if p.Missing(requiredPermissions) {
		return nil, errors.New("I am missing permissions in the channel")
	}

	// Send poll message in channel
	message, err := client.Rest().CreateMessage(channelId, discord.MessageCreate{
		Content:    MakePollText(pollData),
		Embeds:     MakePollEmbeds(pollData),
		Components: MakePollComponents(client, db, pollData),
	})
	if err != nil {
		return nil, err
	}

	// Save poll to database
	pollData.MessageId = message.ID.String()
	if err = CreatePollData(db, pollData); err != nil {
		return nil, err
	}

	return message, nil
}

func CreatePollDM(interaction *events.ApplicationCommandInteractionCreate, db *database.GormDB, data *schema.PollData) (*discord.Message, error) {
	// Send poll message in DM
	err := interaction.CreateMessage(discord.MessageCreate{
		Content:    MakePollText(data),
		Embeds:     MakePollEmbeds(data),
		Components: MakePollComponents(interaction.Client(), db, data),
	})
	if err != nil {
		return nil, err
	}

	// Get interaction response
	message, err := interaction.Client().Rest().GetInteractionResponse(interaction.Client().ID(), interaction.Token())
	if err != nil {
		return nil, err
	}

	// Save poll to database
	data.MessageId = message.ID.String()
	if err = CreatePollData(db, data); err != nil {
		return nil, err
	}

	return message, nil
}

func CreatePollData(db *database.GormDB, pollData *schema.PollData) error {
	return db.Create(pollData).Error
}

func VotePoll(db *database.GormDB, pollData *schema.PollData, userId string, optionIds []int) (string, error) {
	action := ""
	for i, option := range pollData.Options {
		index := slices.Index(option.Voters, userId) // Find index of user's id in option's voters array

		if index != -1 { // They have already voted on the option

			if action == "added" {
				action = "swapped"
			} else {
				action = "removed"
			}

			if pollData.CanSubmit && len(option.Voters) <= 1 {
				db.Delete(&option)
				pollData.Options = append(pollData.Options[:i], pollData.Options[i+1:]...)
			} else {
				pollData.Options[i].Voters = append(option.Voters[:index], option.Voters[index+1:]...)
				if err := db.Save(&pollData.Options[i]).Error; err != nil {
					return "", err
				}
			}

		} else if slices.Index(optionIds, int(option.OptionId)) != -1 { // They haven't voted on the option yet

			if action == "removed" {
				action = "swapped"
			} else {
				action = "added"
			}

			pollData.Options[i].Voters = append(option.Voters, userId)
			if err := db.Save(&pollData.Options[i]).Error; err != nil {
				return "", err
			}

		}
	}
	return action, nil
}

func UpdatePoll(interaction events.InteractionCreate) {
	//interaction.Client().Rest().UpdateInteractionResponse()
}

func EndPoll(client bot.Client, db *database.GormDB, pollData *schema.PollData, enderId *string) error {
	t := time.Now()
	pollData.HasEnded = true
	pollData.EndAt = &t
	pollData.EnderUserId = enderId

	res := db.Save(pollData)
	if res.Error != nil {
		return res.Error
	}

	err := FetchPollUser(client, db, pollData)
	if err != nil {
		return err
	}
	err = FetchPollEnder(client, db, pollData)
	if err != nil {
		return err
	}

	content := MakePollText(pollData)
	pollEmbeds := MakePollEmbeds(pollData)
	pollComponents := MakePollComponents(client, db, pollData)

	message := discord.MessageUpdate{
		Content:    &content,
		Embeds:     &pollEmbeds,
		Components: &pollComponents,
	}

	_, err = client.Rest().UpdateMessage(pollData.ChannelIdSnowflake(), pollData.MessageIdSnowflake(), message)

	return err
}

func EndTimedPollsLoop(client bot.Client, db *database.GormDB, log *zap.Logger) {
	log.Info("Poll loop started")
	for {
		var polls []schema.PollData
		err := db.
			Where("has_ended = FALSE AND end_at < NOW() + INTERVAL '1 minute' AND guild_id != null").
			Order("end_at ASC").
			Preload("Options").
			Find(&polls).
			Error

		if err != nil {
			log.Error("could not fetch polls to end")
			time.Sleep(time.Second * 60)
			continue
		}

		for _, pollData := range polls {
			waitSecs := pollData.EndAt.Unix() - time.Now().Unix()
			if waitSecs > 5 {
				time.Sleep(time.Second * time.Duration(waitSecs))
			}

			err = EndPoll(client, db, &pollData, nil)

			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		waitSecs := 60 - (time.Now().Unix() % 60)
		//fmt.Println("waiting", waitSecs)
		time.Sleep(time.Second * time.Duration(waitSecs))
	}
}

func FetchPollUser(client bot.Client, db *database.GormDB, pollData *schema.PollData) error {
	if pollData.User() != nil {
		return nil
	}

	u, err := db.FetchUser(client, pollData.UserId)
	if err != nil {
		return err
	}
	pollData.SetUser(*u)

	return nil
}

func FetchPollEnder(client bot.Client, db *database.GormDB, pollData *schema.PollData) error {
	if pollData.EnderUserId == nil || pollData.EnderUser() != nil {
		return nil
	}

	u, err := db.FetchUser(client, *pollData.EnderUserId)
	if err != nil {
		return err
	}
	pollData.SetEnderUser(*u)

	return nil
}
