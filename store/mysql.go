package store

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"DVT_Service/conf"
)

// mysqlConnect connect to mysql
func mysqlConnect(cfg *conf.ConfigMysql) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect db")
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect db")
	}
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)
	sqlDB.SetConnMaxLifetime(time.Hour)

	//db.AutoMigrate(
	//	&models.ValidatorInfo{},
	//	&models.ClusterValidatorTask{},
	//	&models.FeeRecipientTask{},
	//	&models.ClusterAmountTask{},
	//	&models.Encrypt{},
	//)

	return db, nil
}
