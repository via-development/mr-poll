package schema

type SuggestionChannel struct {
	Name      string
	GuildId   string
	ChannelId string `gorm:"primaryKey"`

	ApproveChannelId *string
	DenyChannelId    *string

	CreateThreads   bool
	AnonymousAuthor bool
	AnonymousVoters bool
	ShowCounts      bool

	UpvoteEmoji   *string
	DownvoteEmoji *string
	EmbedColor    *int

	Panel SuggestionChannelPanel `gorm:"embedded;embeddedPrefix:panel_"`

	Cooldown int
}

type SuggestionChannelPanel struct {
	Enabled       bool
	Content       string
	Embed         bool
	LastMessageId *string
	Title         *string
	Description   *string
	AuthorName    *string
	AuthorIconUrl *string
}
