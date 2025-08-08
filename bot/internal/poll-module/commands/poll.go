package pollCommands

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/golittie/timeless"
	"github.com/via-development/mr-poll/bot/internal/database"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	pollUtil "github.com/via-development/mr-poll/bot/internal/poll-module/util"
	"github.com/via-development/mr-poll/bot/internal/util"
)

func PollCommand(interaction *events.ApplicationCommandInteractionCreate, db *database.GormDB) error {
	subcommand := interaction.SlashCommandInteractionData().SubCommandGroupName
	if subcommand == nil {
		subcommand = interaction.SlashCommandInteractionData().SubCommandName
	}
	if subcommand == nil {
		return nil
	}

	switch *subcommand {
	case "yes-or-no", "multiple-choice", "single-choice", "submit-choice":
		return pollCreateCommand(interaction, db, *subcommand)
	case "list":
		return pollListCommand(interaction)
	case "end":
		return pollEndCommand(interaction, db)
	case "online":
		return pollOnlineCommand(interaction)
	default:
		return interaction.CreateMessage(discord.MessageCreate{
			Embeds: []discord.Embed{
				util.CommandNotFoundEmbed(),
			},
		})
	}
}

func pollCreateCommand(interaction *events.ApplicationCommandInteractionCreate, db *database.GormDB, subcommand string) error {
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

	// TODO: Timed polls

	pollData := schema.PollData{
		Type:         pollType,
		ChannelId:    interaction.Channel().ID().String(),
		UserId:       interaction.User().ID.String(),
		Question:     question,
		NumOfChoices: 1,
		CanSubmit:    canSubmit,
	}

	if timer != "" {
		endAt := timeless.Parse(timer)
		pollData.EndAt = &endAt
		fmt.Println(timer)
	}

	if pollType == schema.YesOrNoType {
		pollData.Options = []schema.PollOptionData{
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
			opt, exist := cmdData.Option(fmt.Sprint("option-", i+1))
			if !exist {
				break
			}
			name := string(opt.Value)
			name = name[1 : len(name)-1]
			pollData.Options = append(pollData.Options, schema.PollOptionData{
				OptionId: uint(i),
				Name:     name,
				Voters:   []string{},
				Emoji:    util.Alpha[i],
			})
		}
	}

	if pollType == schema.MultipleChoiceType {
		n := cmdData.Int("num-of-choices")
		if n == 0 {
			pollData.NumOfChoices = uint(len(pollData.Options))
		} else {
			pollData.NumOfChoices = uint(n)
		}
	}

	if guild, found := interaction.Guild(); found {
		id := guild.ID.String()
		pollData.GuildId = &id
	}

	err := pollUtil.FetchPollUser(interaction.Client(), db, &pollData)
	if err != nil { // TODO:
		return err
	}

	if pollData.GuildId == nil {
		_, err := pollUtil.CreatePollDM(interaction, db, &pollData)
		return err
	}

	message, err := pollUtil.CreatePoll(interaction.Client(), db, &pollData)

	if err != nil {
		return err
	}

	return interaction.CreateMessage(discord.MessageCreate{
		Flags: discord.MessageFlagEphemeral,
		Embeds: []discord.Embed{
			util.MakeSimpleEmbed(fmt.Sprintf("Created [poll](https://discord.com/channels/%s/%s/%s)!", *pollData.GuildId, message.ChannelID, message.ID)),
		},
	})
}

func pollOnlineCommand(interaction *events.ApplicationCommandInteractionCreate) error {
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

func pollEndCommand(interaction *events.ApplicationCommandInteractionCreate, db *database.GormDB) error {
	cmdData := interaction.SlashCommandInteractionData()

	messageId := cmdData.String("message")

	var pollData schema.PollData
	if err := db.Preload("Options").First(&pollData, messageId).Error; err != nil {
		return interaction.CreateMessage(pollUtil.PollNotFoundMessage())
	}

	if pollData.HasEnded {
		// TODO: poll already ended error
		return nil
	}

	if pollData.GuildId == nil {
		return interaction.CreateMessage(discord.MessageCreate{
			Flags:   discord.MessageFlagEphemeral,
			Content: "Due to a discord limitation, you can only end this poll using the menu button on the poll message!",
		})
	}

	g1 := ""
	if interaction.GuildID() != nil {
		g1 = interaction.GuildID().String()
	}
	g2 := ""
	if pollData.GuildId != nil {
		g2 = *pollData.GuildId
	}
	if g1 != g2 {
		return interaction.CreateMessage(pollUtil.PollNotFoundMessage())
	}

	// If it's not their own poll check for permissions
	if pollData.UserId != interaction.User().ID.String() {
		if interaction.GuildID() == nil {
			// TODO : not your message error
			return interaction.CreateMessage(pollUtil.NotYourPollMessage())
		}

		channel, _ := interaction.Client().Caches().GuildTextChannel(interaction.Channel().ID())
		perms := interaction.Client().Caches().MemberPermissionsInChannel(channel, interaction.Member().Member)

		if !perms.Has(discord.PermissionManageMessages) {
			// TODO : no perms error
			return interaction.CreateMessage(pollUtil.NotYourPollMessage())
		}
	}

	err := pollUtil.FetchPollUser(interaction.Client(), db, &pollData)
	if err != nil {
		return err
	}
	enderId := interaction.User().ID.String()
	return pollUtil.EndPoll(interaction.Client(), db, &pollData, &enderId)
}

func pollListCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	return nil
}
