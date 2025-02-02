package pollUtil

import (
	"github.com/disgoorg/disgo/discord"
	baseUtil "mrpoll_bot/base-util"
	pollDatabase "mrpoll_bot/poll-module/database"
)

// MakePollEmbed makes an embed for a poll with the poll data provided.
func MakePollEmbed(data pollDatabase.PollData) discord.Embed {
	return discord.Embed{
		Title: data.Question,
		Color: baseUtil.Config.EmbedColor,
	}
}
