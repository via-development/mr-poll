package schema

import "time"

type SuggestionCooldown struct {
	CooldownId int    `gorm:"primaryKey"`
	UserId     string `gorm:"not null"`
	ChannelId  string
	CreatedAt  time.Time
}
