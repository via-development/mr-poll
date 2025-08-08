package schema

type UserData struct {
	UserId          string `gorm:"primaryKey"`
	Username        string
	DisplayName     string
	PermissionLevel int

	Stats UserStatsData `gorm:"foreignKey:UserId;references:UserId"`
}

type UserStatsData struct {
	UserId          string `gorm:"primaryKey"`
	PollCount       int
	SuggestionCount int
}
