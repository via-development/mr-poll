package pollUtil

import (
	"github.com/disgoorg/disgo/discord"
	baseUtil "mrpoll_bot/base-util"
)

func MakePollEmbed() discord.Embed {
	return discord.Embed{
		Title: "Poll",
		Color: baseUtil.Config.EmbedColor,
	}
}
