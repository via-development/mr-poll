package database

import (
	"fmt"
	"github.com/gofor-little/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	pollDatabase "mrpoll_bot/poll-module/database"
)

var DB *gorm.DB

func InitDB() {
	dsn, err := env.MustGet("DSN")
	if err != nil {
		panic(err)
	}
	fmt.Println("[Gorm]: Connecting...")
	DB, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		panic(err)
	}
	fmt.Println("[Gorm]: Operational!")
	err = DB.AutoMigrate(&pollDatabase.PollData{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&pollDatabase.PollOptionData{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&pollDatabase.PollRoleData{})
	if err != nil {
		panic(err)
	}
}
