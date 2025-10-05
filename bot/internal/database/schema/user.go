package schema

import "github.com/golittie/timeless/pkg/dateformat"

type User struct {
	UserId          string `gorm:"primaryKey"`
	Username        string
	DisplayName     *string
	PermissionLevel int
	UTCOffset       *float32
	DateFormat      *dateformat.DateFormat
}

func (u *User) SafeName() string {
	if u.DisplayName != nil {
		return *u.DisplayName
	}
	return "@" + u.Username
}
