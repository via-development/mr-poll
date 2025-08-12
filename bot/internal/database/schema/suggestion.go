package schema

import "github.com/lib/pq"

type Suggestion struct {
	ChannelId string
	MessageId string
	UserId    *string

	Title       string
	Description string

	Upvotes   pq.StringArray `gorm:"type:text[]"`
	Downvotes pq.StringArray `gorm:"type:text[]"`
}
