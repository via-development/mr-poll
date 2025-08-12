package schema

type SuggestionChannel struct {
	Name      string
	GuildId   string
	ChannelId string

	ApproveChannelId *string
	DenyChannelId    *string

	CreateThreads   bool
	AnonymousAuthor bool
	AnonymousVoters bool
	ShowCounts      bool

	//upvoteEmoji   String?
	//downvoteEmoji String?
	//embedColor    Int?

	Cooldown int
}
