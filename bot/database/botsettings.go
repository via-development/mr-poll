package database

import "github.com/lib/pq"

type BotSettings struct {
	BotId           string         `gorm:"primaryKey"`
	DisabledModules pq.StringArray `gorm:"type:text[]"`
}
