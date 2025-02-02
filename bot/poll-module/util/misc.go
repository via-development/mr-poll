package pollUtil

import (
	"github.com/disgoorg/disgo/discord"
	pollDatabase "mrpoll_bot/poll-module/database"
)

// PollMessage generates the options for the poll's message, embed and component.
func PollMessage(data pollDatabase.PollData) discord.MessageCreate {
	return discord.MessageCreate{
		Embeds:     []discord.Embed{MakePollEmbed(data)},
		Components: []discord.ContainerComponent{MakePollComponents(data)},
	}
}
