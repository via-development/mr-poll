package database

import (
	"flag"
	"fmt"
	"github.com/gofor-little/env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"slices"
)

var DB *gorm.DB

func InitDB() {
	dsn, err := env.MustGet("DSN")
	if err != nil {
		panic(err)
	}

	databaseDebug := slices.Index(os.Args, "-dd") != -1
	flag.Parse()
	migrateFlag := slices.Index(os.Args, "-m") != -1
	flag.Parse()

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
}
