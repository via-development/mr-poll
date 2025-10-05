package poll_module

import (
	"errors"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/golittie/timeless"
	"github.com/via-development/mr-poll/bot/internal/database"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	"github.com/via-development/mr-poll/bot/internal/util"
	"gorm.io/gorm"
)

func (m *PollModule) PollCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	subcommand := interaction.SlashCommandInteractionData().SubCommandGroupName
	if subcommand == nil {
		subcommand = interaction.SlashCommandInteractionData().SubCommandName
	}
	if subcommand == nil {
		return nil
	}

	switch *subcommand {
	case "yes-or-no", "multiple-choice", "single-choice":
		return m.pollCreateCommand(interaction, m.db, *subcommand)
	case "list":
		return m.pollListCommand(interaction)
	case "end":
		return m.pollEndCommand(interaction, m.db)
	case "online":
		return m.pollOnlineCommand(interaction)
	default:
		return interaction.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{
				util.CommandNotFoundEmbed(),
			},
		})
	}
}

func (m *PollModule) pollCreateCommand(interaction *events.ApplicationCommandInteractionCreate, db *database.GormDB, subcommand string) error {
	// TODO:
	//err := interaction.DeferCreateMessage(false)
	//if err != nil {
	//	return err
	//}
	//user := interaction.User()
	//db.Save(&schema.User{
	//	UserId:      user.ID.String(),
	//	Username:    user.Username,
	//	DisplayName: user.GlobalName,
	//})

	cmdData := interaction.SlashCommandInteractionData()

	question := cmdData.String("question")
	timer := cmdData.String("timer")
	canSubmit := cmdData.Bool("can-submit")

	pollType := schema.YesOrNoType
	switch subcommand {
	case "single-choice":
		pollType = schema.SingleChoiceType
	case "multiple-choice":
		pollType = schema.MultipleChoiceType
	}

	poll := schema.Poll{
		Type:      pollType,
		ChannelId: interaction.Channel().ID().String(),
		UserId:    interaction.User().ID.String(),
		Question:  question,
		CanSubmit: canSubmit,
	}

	if n, f := cmdData.OptInt("num-of-choices"); f {
		n2 := uint(n)
		poll.NumOfChoices = &n2
	}

	err := m.FetchPollUser(&poll)
	if err != nil { // TODO:
		return err
	}

	if timer != "" {
		if poll.User().DateFormat == nil || poll.User().UTCOffset == nil {
			return interaction.CreateMessage(discord.MessageCreate{
				Embeds: []discord.Embed{
					util.MakeErrorEmbed("Please set your timezone and dateformat with `/mytime`!"),
				},
			})
		}

		endAt := timeless.Parse(timer,
			timeless.WithTimezone(*poll.User().UTCOffset),
			timeless.WithDateFormat(*poll.User().DateFormat),
		)

		poll.EndAt = &endAt
	}

	if pollType == schema.YesOrNoType {
		poll.Options = []schema.PollOption{
			{
				OptionId: 0,
				Name:     "Yes",
				Voters:   []string{},
				Emoji:    "#check",
			},
			{
				OptionId: 1,
				Name:     "No",
				Voters:   []string{},
				Emoji:    "#cross",
			},
		}
	} else {
		for i := range 10 {
			name, exist := cmdData.OptString(fmt.Sprint("option-", i+1))
			if !exist {
				break
			}
			poll.Options = append(poll.Options, schema.PollOption{
				OptionId: uint(i),
				Name:     name,
				Voters:   []string{},
				Emoji:    util.Alpha[i],
			})
		}
	}

	if len(poll.Options) == 0 && !poll.CanSubmit {
		return interaction.CreateMessage(discord.MessageCreate{
			Flags: discord.MessageFlagEphemeral,
			Embeds: []discord.Embed{
				util.MakeErrorEmbed("You didn't add any options to the poll!"),
			},
		})
	}

	if guild, f := interaction.Guild(); f {
		id := guild.ID.String()
		poll.GuildId = &id
	}

	if poll.GuildId == nil {
		_, err := m.CreatePollDM(interaction, &poll)
		return err
	}

	message, err := m.CreatePoll(&poll)

	if err != nil {
		return err
	}

	return interaction.CreateMessage(discord.MessageCreate{
		Flags: discord.MessageFlagEphemeral,
		Embeds: []discord.Embed{
			util.MakeSuccessEmbed(fmt.Sprintf("Created [poll](https://discord.com/channels/%s/%s/%s)!", *poll.GuildId, message.ChannelID, message.ID)),
		},
	})
}

func (m *PollModule) pollOnlineCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	var guildId string
	if guild, found := interaction.Guild(); found {
		guildId = guild.ID.String()
	} else {
		guildId = interaction.Channel().ID().String()
	}

	return interaction.CreateMessage(discord.MessageCreate{
		Components: []discord.ContainerComponent{
			discord.ActionRowComponent{
				discord.ButtonComponent{ // Leaked?!?!
					URL:   "https://mrpoll.xyz/polls/" + guildId,
					Style: discord.ButtonStyleLink,
					Label: "View Polls",
				},
			},
		},
	})
}

func (m *PollModule) pollEndCommand(interaction *events.ApplicationCommandInteractionCreate, db *database.GormDB) error {
	cmdData := interaction.SlashCommandInteractionData()

	messageId := cmdData.String("message")

	var guildId *string
	if interaction.GuildID() != nil {
		g := interaction.GuildID().String()
		guildId = &g
	}

	var pollData schema.Poll
	if err := db.Preload("Options").Find(&pollData, schema.Poll{
		MessageId: messageId,
		GuildId:   guildId,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return interaction.CreateMessage(PollNotFoundMessage())
		} else {
			return err
		}
	}

	if pollData.HasEnded {
		// TODO: poll already ended error
		return nil
	}

	if pollData.GuildId == nil || interaction.GuildID() == nil {
		return interaction.CreateMessage(discord.MessageCreate{
			Flags:   discord.MessageFlagEphemeral,
			Content: "Due to Discord limitations, you can only end this poll using the menu button on the poll message!",
		})
	}

	// If it's not their own poll, check for permissions
	if pollData.UserId != interaction.User().ID.String() {
		channel, _ := interaction.Client().Caches().GuildTextChannel(interaction.Channel().ID())
		perms := interaction.Client().Caches().MemberPermissionsInChannel(channel, interaction.Member().Member)

		if !perms.Has(discord.PermissionManageMessages) {
			return interaction.CreateMessage(NotYourPollMessage())
		}
	}

	err := m.FetchPollUser(&pollData)
	if err != nil {
		return err
	}

	enderId := interaction.User().ID.String()

	err = m.EndPoll(&pollData, &enderId)
	if err != nil {
		return err
	}

	err = m.UpdatePollMessage(&pollData)
	if err != nil {
		return err
	}

	return interaction.CreateMessage(discord.MessageCreate{
		Flags: discord.MessageFlagEphemeral,
		Embeds: []discord.Embed{
			util.MakeSuccessEmbed(fmt.Sprintf("Ended [poll](%s)!", pollData.MessageUrl())),
		},
	})
}

func (m *PollModule) pollListCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	return nil
}
