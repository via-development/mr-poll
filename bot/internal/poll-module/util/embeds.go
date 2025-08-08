package pollUtil

import (
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	"github.com/via-development/mr-poll/bot/internal/util"
	"sort"
)

// MakePollEmbeds makes the embeds for a poll with the poll data provided.
func MakePollEmbeds(pollData *schema.PollData) []discord.Embed {
	optionStr := ""
	sort.Slice(pollData.Options, func(i, j int) bool {
		return pollData.Options[i].OptionId < pollData.Options[j].OptionId
	})
	for _, option := range pollData.Options {
		optionStr += fmt.Sprintf("%s `%d votes` %s\n", option.ChatEmoji(), len(option.Voters), option.Name)
	}
	pollEmbeds := []discord.Embed{{
		Author: &discord.EmbedAuthor{
			Name:    pollData.User().DisplayName,
			IconURL: "https://ava.viadev.xyz/" + pollData.UserId,
		},
		Title:       pollData.Question,
		URL:         "https://mrpoll.dev/vote",
		Description: optionStr,
		Color:       util.EmbedColor,
	}}

	if pollData.Images != nil {
		pollEmbeds[0].Image = &discord.EmbedResource{
			URL: (*pollData.Images)[0],
		}

		for i := range len(*pollData.Images) - 1 {
			pollEmbeds = append(pollEmbeds, discord.Embed{
				URL: "https://mrpoll.dev/vote",
				Image: &discord.EmbedResource{
					URL: (*pollData.Images)[i],
				},
			})
		}

	}

	return pollEmbeds
}
