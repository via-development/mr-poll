package schema

type User struct {
	UserId          string `gorm:"primaryKey"`
	Username        string
	DisplayName     string
	PermissionLevel int

	Stats UserStats `gorm:"foreignKey:UserId;references:UserId"`
}

type UserStats struct {
	UserId          string `gorm:"primaryKey"`
	PollCount       int
	SuggestionCount int
}
