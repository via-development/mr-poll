package schema

type SuggestionChannel struct {
	Name      string
	GuildId   string
	ChannelId string

	ApproveChannelId *string
	DenyChannelId    *string

	createThreads   bool
	anonymousAuthor bool
	anonymousVoters bool
	showCounts      bool

	//upvoteEmoji   String?
	//downvoteEmoji String?
	//embedColor    Int?
	cooldown int
}
