package schema

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/lib/pq"
	"github.com/via-development/mr-poll/bot/internal/util"
	"time"
)

type Poll struct {
	Type      PollType `gorm:"not null"`
	MessageId string   `gorm:"primaryKey"`
	ChannelId string   `gorm:"not null"`
	GuildId   *string

	UserId string `gorm:"not null"`
	user   *User

	Question      string          `gorm:"not null"`
	Options       []PollOption    `gorm:"foreignKey:MessageId;references:MessageId"`
	PollRoles     []PollRole      `gorm:"foreignKey:MessageId;references:MessageId"`
	AnonymousMode AnonymousType   `gorm:"not null"`
	CanSubmit     bool            `gorm:"default:false"`
	NumOfChoices  uint            `gorm:"not null"`
	Images        *pq.StringArray `gorm:"type:text[]"`

	HasEnded    bool `gorm:"default:false"`
	EndAt       *time.Time
	EnderUserId *string
	enderUser   *User
}

type PollType = uint

const (
	YesOrNoType PollType = iota
	SingleChoiceType
	MultipleChoiceType
	RatingType
	PetitionType
)

type AnonymousType = uint

const (
	AnonymousNone AnonymousType = iota
	AnonymousForever
	AnonymousUntilEnd
)

type PollOption struct {
	Uid       uint           `gorm:"primaryKey"`
	OptionId  uint           `gorm:"not null"`
	MessageId string         `gorm:"not null"`
	Name      string         `gorm:"not null"`
	Emoji     string         `gorm:"not null"`
	Voters    pq.StringArray `gorm:"type:text[]"`
	SubmitBy  *string
}

type PollRole struct {
	Uid        uint `gorm:"primaryKey"`
	MessageId  string
	RoleId     string
	BonusCount uint
	Required   bool
}

func (p *Poll) MessageIdSnowflake() snowflake.ID {
	return snowflake.MustParse(p.MessageId)
}

func (p *Poll) ChannelIdSnowflake() snowflake.ID {
	return snowflake.MustParse(p.ChannelId)
}

func (p *Poll) GuildIdSnowflake() *snowflake.ID {
	if p.GuildId == nil {
		return nil
	}
	s := snowflake.MustParse(*p.GuildId)
	return &s
}

func (p *Poll) UserIdSnowflake() *snowflake.ID {
	s := snowflake.MustParse(p.UserId)
	return &s
}

func (p *Poll) SetUser(user User) {
	p.user = &user
}

func (p *Poll) User() *User {
	return p.user
}

func (p *Poll) SetEnderUser(user User) {
	p.enderUser = &user
}

func (p *Poll) EnderUser() *User {
	return p.enderUser
}

func (o *PollOption) parseEmoji() (string, bool) {
	emoji := o.Emoji
	if emoji[0] == '#' {
		switch emoji[1:] {
		case "check":
			emoji = "1268234822304792676"
		case "cross":
			emoji = "1268234748988493905"
		default:
			panic(emoji + " is not an emoji")
		}
	}
	isUnicode := !util.NumRegex.Match([]byte(emoji))
	return emoji, isUnicode
}

func (o *PollOption) ApiEmoji() discord.ComponentEmoji {
	emoji, isUnicode := o.parseEmoji()

	if isUnicode {
		return discord.ComponentEmoji{
			Name: emoji,
		}
	}

	return discord.ComponentEmoji{
		ID: snowflake.MustParse(emoji), Name: "e",
	}
}

func (o *PollOption) ChatEmoji() string {
	emoji, isUnicode := o.parseEmoji()

	if isUnicode {
		return emoji
	}

	return "<:e:" + emoji + ">"
}
