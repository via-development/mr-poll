package pollButtons

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"mrpoll_bot/database"
	pollUtil "mrpoll_bot/poll-module/util"
	"strings"
)

func PollMenuButton(interaction *events.ComponentInteractionCreate) error {
	var pollData database.PollData
	res := database.DB.Preload("Options").First(&pollData, interaction.Message.ID.String())
	if res.Error != nil {
		return interaction.CreateMessage(pollUtil.PollNotFoundMessage())
	}
	optStr := ""
	for _, opt := range pollData.Options {
		optStr += fmt.Sprintf("%s - %s\n", opt.Name, strings.Join(opt.Voters, ", "))
	}
	return interaction.CreateMessage(discord.MessageCreate{
		Flags:   discord.MessageFlagEphemeral,
		Content: fmt.Sprintf("ℹ️ You can do the same actions by right clicking the poll!\n> %s\nOptions:\n%s", pollData.Question, optStr),
	})
}
