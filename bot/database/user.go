package database

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/snowflake/v2"
)

type UserData struct {
	UserId          string `gorm:"primaryKey"`
	Username        string
	DisplayName     string
	PermissionLevel int

	Stats UserStatsData `gorm:"foreignKey:UserId;references:UserId"`
}

type UserStatsData struct {
	UserId          string `gorm:"primaryKey"`
	PollCount       int
	SuggestionCount int
}

func FetchUser(userId string, client bot.Client) *UserData {
	var userData *UserData
	if DB.First(&userData, userId).Error != nil {
		user, err := client.Rest().GetUser(snowflake.MustParse(userId))
		if err != nil {
			return nil
		}
		return &UserData{
			UserId:      userId,
			Username:    user.Username,
			DisplayName: *user.GlobalName,
			Stats: UserStatsData{
				PollCount:       0,
				SuggestionCount: 0,
			},
		}
	}
	return userData
}
