package pollDatabase

import (
	"github.com/disgoorg/snowflake/v2"
	"strconv"
)

type PollOptionData struct {
	OptionId  uint `gorm:"primaryKey"`
	MessageId string
	Name      string
	Emoji     string
	Voters    []string `gorm:"type:text"`
}

type PollRoleData struct {
	PollRoleId uint `gorm:"primaryKey"`
	MessageId  string
	RoleId     string
	BonusCount uint
	Required   bool
}

type PollData struct {
	Type          uint
	MessageId     string `gorm:"primaryKey"`
	ChannelId     string
	GuildId       string
	Question      string
	AnonymousMode uint
	Options       []PollOptionData `gorm:"foreignKey:MessageId;references:MessageId"`
	PollRoles     []PollRoleData   `gorm:"foreignKey:MessageId;references:MessageId"`
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

func (p *PollData) MessageIdSnowflake() snowflake.ID {
	n, _ := strconv.Atoi(p.MessageId)
	return snowflake.ID(n)
}

func (p *PollData) ChannelIdSnowflake() snowflake.ID {
	n, _ := strconv.Atoi(p.ChannelId)
	return snowflake.ID(n)
}

func (p *PollData) GuildIdSnowflake() snowflake.ID {
	n, _ := strconv.Atoi(p.GuildId)
	return snowflake.ID(n)
}
