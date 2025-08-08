package schema

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/snowflake/v2"
	"github.com/lib/pq"
	"github.com/via-development/mr-poll/bot/internal/util"
	"time"
)

type PollData struct {
	Type      PollType `gorm:"not null"`
	MessageId string   `gorm:"primaryKey"`
	ChannelId string   `gorm:"not null"`
	GuildId   *string

	UserId string `gorm:"not null"`
	user   *UserData

	Question      string           `gorm:"not null"`
	Options       []PollOptionData `gorm:"foreignKey:MessageId;references:MessageId"`
	PollRoles     []PollRoleData   `gorm:"foreignKey:MessageId;references:MessageId"`
	AnonymousMode AnonymousType    `gorm:"not null"`
	CanSubmit     bool             `gorm:"default:false"`
	NumOfChoices  uint             `gorm:"not null"`
	Images        *pq.StringArray  `gorm:"type:text[]"`

	HasEnded    bool `gorm:"default:false"`
	EndAt       *time.Time
	EnderUserId *string
	enderUser   *UserData
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

type PollOptionData struct {
	Uid       uint           `gorm:"primaryKey"`
	OptionId  uint           `gorm:"not null"`
	MessageId string         `gorm:"not null"`
	Name      string         `gorm:"not null"`
	Emoji     string         `gorm:"not null"`
	Voters    pq.StringArray `gorm:"type:text[]"`
	SubmitBy  *string
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

func (p *PollData) UserIdSnowflake() *snowflake.ID {
	s := snowflake.MustParse(p.UserId)
	return &s
}

func (p *PollData) SetUser(user UserData) {
	p.user = &user
}

func (p *PollData) User() *UserData {
	return p.user
}

func (p *PollData) SetEnderUser(user UserData) {
	p.enderUser = &user
}

func (p *PollData) EnderUser() *UserData {
	return p.enderUser
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
