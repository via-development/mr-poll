package schema

type UserStats struct {
	UserId          string `gorm:"primaryKey"`
	PollCount       int
	SuggestionCount int
}
