package pollUtil

import (
	"encoding/json"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"mrpoll_bot/database"
	"mrpoll_bot/util"
	"sort"
)

// MakePollEmbeds makes the embeds for a poll with the poll data provided.
func MakePollEmbeds(data database.PollData) []discord.Embed {
	optionStr := ""
	sort.Slice(data.Options, func(i, j int) bool {
		return data.Options[i].OptionId < data.Options[j].OptionId
	})
	for _, option := range data.Options {
		optionStr += fmt.Sprintf("%s `%d votes` %s\n", option.ChatEmoji(), len(option.Voters), option.Name)
	}
	ud, _ := json.Marshal(data.User())
	fmt.Printf("%v\n", string(ud))
	fmt.Printf("%v\n", discord.EmbedAuthor{
		Name:    data.User().DisplayName,
		IconURL: "https://ava.viadev.xyz/" + data.UserId,
	})
	pollEmbeds := []discord.Embed{{
		Author: &discord.EmbedAuthor{
			Name:    data.User().DisplayName,
			IconURL: "https://ava.viadev.xyz/" + data.UserId,
		},
		Title:       data.Question,
		URL:         "https://mrpoll.xyz/vote",
		Description: optionStr,
		Color:       util.Config.EmbedColor,
	}}

	if data.Images != nil {
		pollEmbeds[0].Image = &discord.EmbedResource{
			URL: (*data.Images)[0],
		}

		for i := range len(*data.Images) - 1 {
			pollEmbeds = append(pollEmbeds, discord.Embed{
				URL: "https://mrpoll.xyz/vote",
				Image: &discord.EmbedResource{
					URL: (*data.Images)[i],
				},
			})
		}

	}

	return pollEmbeds
}
