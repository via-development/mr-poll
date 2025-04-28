package database

import (
	"fmt"
	"github.com/gofor-little/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"slices"
)

var DB *gorm.DB
var BotSettingsC BotSettings

func InitDB() {
	dsn, err := env.MustGet("DSN")
	if err != nil {
		panic(err)
	}

	databaseDebug := slices.Index(os.Args, "-dd") != -1
	migrateFlag := slices.Index(os.Args, "-m") != -1

	l := logger.Silent
	if databaseDebug {
		l = logger.Info
	}

	fmt.Println("[Gorm]: Connecting...")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(l),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("[Gorm]: Operational!")

	if migrateFlag {
		migrateSchemas()
		os.Exit(0)
	}

	var botId string
	if botId, err = env.MustGet("BOT_ID"); err != nil {
		panic("BOT_ID environment variable not set")
	}

	DB.First(&BotSettingsC, &BotSettings{BotId: botId})
	if BotSettingsC.BotId == "" {
		BotSettingsC = BotSettings{
			BotId:           botId,
			DisabledModules: []string{},
		}
		err = DB.Create(&BotSettingsC).Error
		if err != nil {
			panic(err)
		}
	}
}

func migrateSchemas() {
	err := DB.AutoMigrate(&PollData{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&PollOptionData{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&PollRoleData{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&BotSettings{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&UserData{})
	if err != nil {
		panic(err)
	}
	err = DB.AutoMigrate(&UserStatsData{})
	if err != nil {
		panic(err)
	}
}
