package database

import (
	"errors"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/snowflake/v2"
	"github.com/via-development/mr-poll/bot/internal/config"
	"github.com/via-development/mr-poll/bot/internal/database/schema"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB

	config *config.Config
	log    *zap.Logger

	botSettings *schema.BotSettings
}

func New(lc fx.Lifecycle, config *config.Config, log *zap.Logger) (*Database, error) {
	db := &Database{
		config: config,
		log:    log,
	}

	log.Info("connecting to the database")

	var err error
	db.DB, err = gorm.Open(postgres.Open(config.DSN))
	if err != nil {
		return nil, err
	}

	log.Info("connected to the database")

	//databaseDebug := slices.Index(os.Args, "-dd") != -1
	//migrateFlag := slices.Index(os.Args, "-m") != -1

	//l := logger.Silent
	//if databaseDebug {
	//	l = logger.Info
	//}

	return db, nil
}

func (db *Database) BotSettings() schema.BotSettings {
	var botSettings schema.BotSettings
	res := db.Where(db.config.BotId).Find(&botSettings)

	if res.RowsAffected == 0 {
		botSettings = schema.BotSettings{
			BotId:           db.config.BotId,
			DisabledModules: []string{},
		}
		db.Create(&botSettings)
	}

	return botSettings
}

func (db *Database) RunMigrations() error {
	return db.AutoMigrate(
		&schema.Poll{},
		&schema.PollOption{},
		&schema.PollRole{},
		&schema.BotSettings{},
		&schema.User{},
		&schema.UserStats{},
		&schema.Suggestion{},
		&schema.SuggestionChannel{},
	)
}

func (db *Database) FetchUser(client *bot.Client, userId string) (*schema.User, error) {
	var userData *schema.User
	err := db.Find(&userData, userId).Error

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) || userData.Username == "" {
		user, err := client.Rest.GetUser(snowflake.MustParse(userId))
		if err != nil {
			return nil, err
		}

		userData = &schema.User{
			UserId:      userId,
			Username:    user.Username,
			DisplayName: user.GlobalName,
		}
		err = db.Create(&userData).Error
		if err != nil {
			db.log.Error("could not save user data", zap.Error(err))
		}
	} else if err != nil {
		return nil, err
	}

	return userData, nil
}
