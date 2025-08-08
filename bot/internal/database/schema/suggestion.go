package schema

import "github.com/lib/pq"

type Suggestion struct {
	ChannelId string
	MessageId string
	UserId    *string

	Title       string
	Description string

	Upvotes   pq.StringArray
	Downvotes pq.StringArray
}
