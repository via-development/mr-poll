package schema

import (
	"github.com/disgoorg/snowflake/v2"
	"github.com/lib/pq"
)

type Suggestion struct {
	GuildId   string
	ChannelId string
	MessageId string `gorm:"primaryKey"`
	UserId    string
	user      *User

	Title       *string
	Description string

	Upvotes   pq.StringArray `gorm:"type:text[]"`
	Downvotes pq.StringArray `gorm:"type:text[]"`

	AnonymousAuthor bool `gorm:"default:false"`
	AnonymousVoters bool `gorm:"default:false"`
	ShowCounts      bool `gorm:"default:false"`
}

func (s *Suggestion) MessageIdSnowflake() snowflake.ID {
	return snowflake.MustParse(s.MessageId)
}

func (s *Suggestion) ChannelIdSnowflake() snowflake.ID {
	return snowflake.MustParse(s.ChannelId)
}

func (s *Suggestion) GuildIdSnowflake() snowflake.ID {
	return snowflake.MustParse(s.GuildId)
}

func (s *Suggestion) UserIdSnowflake() snowflake.ID {
	return snowflake.MustParse(s.UserId)
}

func (s *Suggestion) SetUser(user User) {
	s.user = &user
}

func (s *Suggestion) User() *User {
	return s.user
}
