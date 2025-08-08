package schema

import "time"

type StatsData struct {
	Date            time.Time `gorm:"primaryKey"`
	PollCount       uint
	SuggestionCount uint
}
