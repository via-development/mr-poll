package database

import (
	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

type GormLogger struct {
	logger.Interface
	zap *zap.Logger
}

func LogLevel() {

}
