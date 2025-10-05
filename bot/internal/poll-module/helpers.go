package poll_module

import (
	"errors"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/labstack/gommon/log"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	"github.com/via-development/mr-poll/bot/internal/util"
	"go.uber.org/zap"
	"math"
	"slices"
	"sort"
	"time"
)

var requiredPermissions = discord.PermissionSendMessages // | discord.PermissionViewChannel

// CreatePoll is used by both the website and poll command for creating polls, it checks permissions, creates poll message and saves to database.
func (m *PollModule) CreatePoll(poll *schema.Poll) (*discord.Message, error) {
	channelId := poll.ChannelIdSnowflake()

	// Channel permission check
	channel, found := m.client.Caches().Channel(channelId)
	if !found {
		return nil, errors.New("channel is not in cache")
	}
	member, found := m.client.Caches().Member(*poll.GuildIdSnowflake(), m.client.ApplicationID())
	if !found {
		m, err := m.client.Rest().GetMember(*poll.GuildIdSnowflake(), m.client.ApplicationID())
		if err != nil {
			return nil, err
		}
		member = *m
	}
	p := m.client.Caches().MemberPermissionsInChannel(channel, member)
	if p.Missing(requiredPermissions) {
		return nil, errors.New("I am missing permissions in the channel")
	}

	// Send poll message in channel
	message, err := m.client.Rest().CreateMessage(channelId, discord.MessageCreate{
		Content:    MakePollText(poll),
		Embeds:     m.MakePollEmbeds(poll),
		Components: m.MakePollComponents(poll),
	})
	if err != nil {
		return nil, err
	}

	// Save poll to database
	poll.MessageId = message.ID.String()
	if err = m.CreatePollData(poll); err != nil {
		return nil, err
	}

	return message, nil
}

func (m *PollModule) CreatePollDM(interaction *events.ApplicationCommandInteractionCreate, poll *schema.Poll) (*discord.Message, error) {
	// Send poll message in DM
	err := interaction.CreateMessage(discord.MessageCreate{
		Content:    MakePollText(poll),
		Embeds:     m.MakePollEmbeds(poll),
		Components: m.MakePollComponents(poll),
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
	poll.MessageId = message.ID.String()
	if err = m.CreatePollData(poll); err != nil {
		return nil, err
	}

	return message, nil
}

func (m *PollModule) CreatePollData(poll *schema.Poll) error {
	return m.db.Create(poll).Error
}

func (m *PollModule) VotePoll(poll *schema.Poll, userId string, optionIds []int) (string, error) {
	action := ""
	for i, option := range poll.Options {
		index := slices.Index(option.Voters, userId) // Find index of user's id in option's voters array

		if index != -1 { // They have already voted on the option

			if action == "added" {
				action = "swapped"
			} else {
				action = "removed"
			}

			if poll.CanSubmit && len(option.Voters) <= 1 {
				m.db.Delete(&option)
				poll.Options = append(poll.Options[:i], poll.Options[i+1:]...)
			} else {
				poll.Options[i].Voters = append(option.Voters[:index], option.Voters[index+1:]...)
				if err := m.db.Save(&poll.Options[i]).Error; err != nil {
					return "", err
				}
			}

		} else if slices.Index(optionIds, int(option.OptionId)) != -1 { // They haven't voted on the option yet

			if action == "removed" {
				action = "swapped"
			} else {
				action = "added"
			}

			poll.Options[i].Voters = append(option.Voters, userId)
			if err := m.db.Save(&poll.Options[i]).Error; err != nil {
				return "", err
			}

		}
	}
	return action, nil
}

func (m *PollModule) UpdatePollMessage(poll *schema.Poll) error {
	content := MakePollText(poll)
	pollEmbeds := m.MakePollEmbeds(poll)
	pollComponents := m.MakePollComponents(poll)

	message := discord.MessageUpdate{
		Content:    &content,
		Embeds:     &pollEmbeds,
		Components: &pollComponents,
	}

	_, err := m.client.Rest().UpdateMessage(poll.ChannelIdSnowflake(), poll.MessageIdSnowflake(), message)
	return err
}

func (m *PollModule) EndPoll(poll *schema.Poll, enderId *string) error {
	t := time.Now()
	poll.HasEnded = true
	poll.EndAt = &t
	poll.EnderUserId = enderId

	res := m.db.Save(poll)
	if res.Error != nil {
		return res.Error
	}

	err := m.FetchPollUser(poll)
	if err != nil {
		return err
	}

	return m.FetchPollEnder(poll)
}

func (m *PollModule) EndTimedPollsLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	m.log.Info("poll loop started")
	for {
		go m.endTimedPolls()
		<-ticker.C
	}
}

func (m *PollModule) endTimedPolls() {
	var polls []schema.Poll
	res := m.db.
		Where("has_ended = FALSE AND end_at < (NOW() + INTERVAL '1 minute') AND guild_id IS NOT NULL").
		Order("end_at ASC").
		Preload("Options").
		Find(&polls)

	if res.Error != nil {
		log.Error("could not fetch polls to end")
		return
	}

	if res.RowsAffected == 0 {
		log.Info("no polls to end")
		return
	}

	for _, pollData := range polls {
		waitSecs := pollData.EndAt.Unix() - time.Now().Unix()
		if waitSecs > 5 {
			time.Sleep(time.Second * time.Duration(waitSecs))
		}

		err := m.EndPoll(&pollData, nil)
		if err != nil {
			log.Error("could not end timed poll", zap.Error(err))
			continue
		}

		err = m.UpdatePollMessage(&pollData)
		if err != nil {
			log.Error("could not end timed poll", zap.Error(err))
			continue
		}
	}
}

func (m *PollModule) FetchPollUser(poll *schema.Poll) error {
	if poll.User() != nil {
		return nil
	}

	u, err := m.db.FetchUser(m.client, poll.UserId)
	if err != nil {
		return err
	}
	poll.SetUser(*u)

	return nil
}

func (m *PollModule) FetchPollEnder(pollData *schema.Poll) error {
	if pollData.EnderUserId == nil || pollData.EnderUser() != nil {
		return nil
	}

	u, err := m.db.FetchUser(m.client, *pollData.EnderUserId)
	if err != nil {
		return err
	}
	pollData.SetEnderUser(*u)

	return nil
}

// MakePollEmbeds makes the embeds for a poll with the poll data provided.
func (m *PollModule) MakePollEmbeds(pollData *schema.Poll) []discord.Embed {
	optionStr := ""
	sort.Slice(pollData.Options, func(i, j int) bool {
		return pollData.Options[i].OptionId < pollData.Options[j].OptionId
	})
	for _, option := range pollData.Options {
		optionStr += fmt.Sprintf("%s `%d votes` %s\n", option.ChatEmoji(), len(option.Voters), option.Name)
	}
	if optionStr == "" {
		optionStr = "No options submitted yet!"
	}
	pollEmbeds := []discord.Embed{{
		Author: &discord.EmbedAuthor{
			Name:    pollData.User().SafeName() + " asked",
			IconURL: "https://ava.viadev.xyz/" + pollData.UserId,
		},
		Title:       pollData.Question,
		URL:         "https://mrpoll.dev/vote",
		Description: optionStr,
		Color:       util.EmbedColor,
	}}

	if !pollData.HasEnded {
		footerText := fmt.Sprintf("@%s • Vote for ", pollData.User().Username)
		if pollData.CanSubmit {
			footerText += "or submit "
		}
		footerText += "your option!"
		pollEmbeds[0].Footer = &discord.EmbedFooter{
			Text: footerText,
		}
	}

	if pollData.Images != nil {
		pollEmbeds[0].Image = &discord.EmbedResource{
			URL: (*pollData.Images)[0],
		}

		for i := range len(*pollData.Images) - 1 {
			pollEmbeds = append(pollEmbeds, discord.Embed{
				URL: util.BotVoteURL,
				Image: &discord.EmbedResource{
					URL: (*pollData.Images)[i],
				},
			})
		}

	}

	return pollEmbeds
}

var menuButton = discord.ButtonComponent{
	Emoji:    &discord.ComponentEmoji{ID: util.OptionsEmoji, Name: "e"},
	CustomID: "poll:menu",
	Style:    discord.ButtonStyleSecondary,
}

// MakePollComponents makes components for a poll with the poll data provided.
func (m *PollModule) MakePollComponents(data *schema.Poll) []discord.ContainerComponent {
	var components []discord.ContainerComponent
	switch data.Type {
	case schema.YesOrNoType, schema.SingleChoiceType:
		var options discord.ActionRowComponent
		for i, op := range data.Options {
			e := op.ApiEmoji()
			s := discord.ButtonStylePrimary
			if data.Type == schema.YesOrNoType {
				if i == 0 {
					s = discord.ButtonStyleSuccess
				} else {
					s = discord.ButtonStyleDanger
				}
			}
			options = append(options, discord.ButtonComponent{
				Emoji:    &e,
				CustomID: fmt.Sprint("poll:option-", i),
				Style:    s,
				Disabled: data.HasEnded,
			})
		}
		if len(options) < 10 && data.CanSubmit {
			options = append(options, discord.ButtonComponent{
				CustomID: "poll:option-submit",
				Style:    discord.ButtonStyleSecondary,
				Emoji:    &discord.ComponentEmoji{ID: util.SubmitEmoji, Name: "e"},
			})
		}
		options = append(options, menuButton)
		for i := range (len(options) / 5) + 1 {
			upper := int(math.Min(float64(i*5+5), float64(len(options))))
			components = append(components, options[i*5:upper])
		}
	case schema.MultipleChoiceType:
		var options []discord.StringSelectMenuOption
		for _, opt := range data.Options {
			e := opt.ApiEmoji()
			desc := ""
			if data.CanSubmit && opt.SubmitBy != nil {
				user, err := m.db.FetchUser(m.client, *opt.SubmitBy)
				if err == nil && user != nil {
					desc = fmt.Sprint("Submitted by @", user.Username)
				}
			}
			options = append(options, discord.StringSelectMenuOption{
				Label:       opt.Name,
				Value:       fmt.Sprint("option-", opt.OptionId),
				Emoji:       &e,
				Description: desc,
			})
		}
		if len(options) < 10 && data.CanSubmit {
			options = append(options, discord.StringSelectMenuOption{
				Label: "Submit your answer!",
				Value: "option:submit",
			})
		}

		selectmenu := discord.StringSelectMenuComponent{
			CustomID: "poll:opts",
			Options:  options,
			Disabled: data.HasEnded,
		}
		if data.NumOfChoices != nil {
			selectmenu.MaxValues = int(*data.NumOfChoices)
		}

		components = []discord.ContainerComponent{
			discord.ActionRowComponent{selectmenu},
			discord.ActionRowComponent{menuButton},
		}
	default:
		components[0] = discord.ActionRowComponent{
			discord.ButtonComponent{
				Label:    "Something went wrong!",
				CustomID: "oops",
				Style:    discord.ButtonStyleSecondary,
				Disabled: true,
			},
		}
	}

	return components
}

func MakePollText(poll *schema.Poll) string {
	if poll.HasEnded {
		if poll.EnderUserId == nil {
			return fmt.Sprintf("🛑 This poll was ended automatically. (<t:%d:R>)", poll.EndAt.Unix())
		}
		return fmt.Sprintf("🛑 This poll was ended by <@%s> (@%s). (<t:%d:R>)", *poll.EnderUserId, poll.EnderUser().Username, poll.EndAt.Unix())
	}
	if poll.EndAt != nil {
		if poll.GuildId == nil {
			return "⏱️ This poll will not end automatically, this is a dm. How did you manage to do this?"
		}
		return fmt.Sprintf("⏱️ This poll will end automatically. (<t:%d:R>)", poll.EndAt.Unix())
	}
	return ""
}

func PollNotFoundMessage() discord.MessageCreate {
	return discord.MessageCreate{
		Flags:   discord.MessageFlagEphemeral,
		Content: "Cannot fetch this poll!",
	}
}

func NotYourPollMessage() discord.MessageCreate {
	return discord.MessageCreate{
		Flags:   discord.MessageFlagEphemeral,
		Content: "This isn't your poll!",
	}
}

func PollOptionSubmitModel(messageId string) discord.ModalCreate {
	return discord.ModalCreate{
		Title:    "Submit your answer!",
		CustomID: "poll:option-submit:" + messageId,
		Components: []discord.ContainerComponent{
			discord.ActionRowComponent{
				discord.TextInputComponent{
					Label:    "Your Answer",
					CustomID: "answer",
					Style:    discord.TextInputStyleShort,
				},
			},
		},
	}
}
