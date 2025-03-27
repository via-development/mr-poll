package pollUtil

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"mrpoll_bot/database"
	"mrpoll_bot/util"
	"sort"
)

// MakePollEmbed makes an embed for a poll with the poll data provided.
func MakePollEmbed(data database.PollData) discord.Embed {
	optionStr := ""
	sort.Slice(data.Options, func(i, j int) bool {
		return data.Options[i].OptionId < data.Options[j].OptionId
	})
	for _, option := range data.Options {
		optionStr += fmt.Sprintf("%s `%d votes` %s\n", option.ChatEmoji(), len(option.Voters), option.Name)
	}
	return discord.Embed{
		Author: &discord.EmbedAuthor{
			Name:    "Someone asked",
			IconURL: "https://ava.viadev.xyz/" + data.UserId,
		},
		Title:       data.Question,
		Description: optionStr,
		Color:       util.Config.EmbedColor,
	}
}
