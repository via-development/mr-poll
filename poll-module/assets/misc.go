package pollAssets

const (
	FlagPollEnded               = 1 << iota
	FlagPollAnonymousVoting     = 1 << iota
	FlagPollDontShowCount       = 1 << iota
	FlagPollEmojiOptions        = 1 << iota
	FlagPollBonusVoteRolesStack = 1 << iota
)
