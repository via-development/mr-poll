package database

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/lib/pq"
	"mrpoll_bot/util"
	"time"
)

type PollData struct {
	Type      uint   `gorm:"not null"`
	MessageId string `gorm:"primaryKey"`
	ChannelId string `gorm:"not null"`
	GuildId   *string
	UserId    string `gorm:"not null"`

	Question      string           `gorm:"not null"`
	Options       []PollOptionData `gorm:"foreignKey:MessageId;references:MessageId"`
	PollRoles     []PollRoleData   `gorm:"foreignKey:MessageId;references:MessageId"`
	AnonymousMode uint             `gorm:"not null"`
	NumOfChoices  uint             `gorm:"not null"`

	EndedAt     *time.Time
	EnderUserId *string
}

type PollType = uint

const (
	YesOrNoType PollType = iota
	SingleChoiceType
	MultipleChoiceType
	SubmitChoiceType
)

type AnonymousType = uint

const (
	AnonymousNone AnonymousType = iota
	AnonymousForever
	AnonymousUntilEnd
)

type PollOptionData struct {
	Uid       uint           `gorm:"primaryKey"`
	OptionId  uint           `gorm:"not null"`
	MessageId string         `gorm:"not null"`
	Name      string         `gorm:"not null"`
	Emoji     string         `gorm:"not null"`
	Voters    pq.StringArray `gorm:"type:text[]"`
}

type PollRoleData struct {
	Uid        uint `gorm:"primaryKey"`
	MessageId  string
	RoleId     string
	BonusCount uint
	Required   bool
}

func (p *PollData) MessageIdSnowflake() snowflake.ID {
	return snowflake.MustParse(p.MessageId)
}

func (p *PollData) ChannelIdSnowflake() snowflake.ID {
	return snowflake.MustParse(p.ChannelId)
}

func (p *PollData) GuildIdSnowflake() *snowflake.ID {
	if p.GuildId == nil {
		return nil
	}
	s := snowflake.MustParse(*p.GuildId)
	return &s
}

func (o *PollOptionData) parseEmoji() (string, bool) {
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

func (o *PollOptionData) ApiEmoji() discord.ComponentEmoji {
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

func (o *PollOptionData) ChatEmoji() string {
	emoji, isUnicode := o.parseEmoji()

	if isUnicode {
		return emoji
	}

	return "<:e:" + emoji + ">"
}
