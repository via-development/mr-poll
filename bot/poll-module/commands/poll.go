package pollCommands

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"mrpoll_bot/database"
	pollUtil "mrpoll_bot/poll-module/util"
	"mrpoll_bot/util"
)

func PollCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	subcommand := interaction.SlashCommandInteractionData().SubCommandGroupName
	if subcommand == nil {
		subcommand = interaction.SlashCommandInteractionData().SubCommandName
	}
	if subcommand == nil {
		return nil
	}

	switch *subcommand {
	case "yes-or-no", "multiple-choice", "single-choice", "submit-choice":
		return pollCreateCommand(interaction, *subcommand)
	case "list":
		return pollListCommand(interaction)
	case "end":
		return pollEndCommand(interaction)
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

func pollCreateCommand(interaction *events.ApplicationCommandInteractionCreate, subcommand string) error {
	cmdData := interaction.SlashCommandInteractionData()

	question := cmdData.String("question")
	timer := cmdData.String("timer")
	fmt.Println(timer)

	pollType := database.YesOrNoType
	switch subcommand {
	case "single-choice":
		pollType = database.SingleChoiceType
	case "multiple-choice":
		pollType = database.MultipleChoiceType
	case "submit-choice":
		pollType = database.SubmitChoiceType
	}

	pollData := database.PollData{
		Type:         pollType,
		ChannelId:    interaction.Channel().ID().String(),
		UserId:       interaction.User().ID.String(),
		Question:     question,
		NumOfChoices: 1,
	}

	if pollType == database.YesOrNoType {
		pollData.Options = []database.PollOptionData{
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
	} else if pollType != database.SubmitChoiceType {
		for i := range 10 {
			opt, exist := cmdData.Option(fmt.Sprint("option-", i+1))
			if !exist {
				break
			}
			name := string(opt.Value)
			name = name[1 : len(name)-1]
			pollData.Options = append(pollData.Options, database.PollOptionData{
				OptionId: uint(i),
				Name:     name,
				Voters:   []string{},
				Emoji:    util.Alpha[i],
			})
		}
	}

	if pollType == database.MultipleChoiceType {
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
	pollData.FetchUser(interaction.Client())

	if pollData.GuildId == nil {
		_, err := pollUtil.CreatePollDM(interaction, pollData)
		return err
	}

	message, err := pollUtil.CreatePoll(interaction.Client(), pollData)

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
					URL:   "https://mrpoll.dev/polls/" + guildId,
					Style: discord.ButtonStyleLink,
					Label: "View Polls",
				},
			},
		},
	})
}

func pollEndCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	return nil
}

func pollListCommand(interaction *events.ApplicationCommandInteractionCreate) error {
	return nil
}
