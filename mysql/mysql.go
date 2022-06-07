package mysql

import (
	"errors"
	"fmt"
	"github.com/wonderivan/logger"
	gmsyql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"lubanKubernets/config"
	"lubanKubernets/model"
	"time"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DbUser, config.DbPwd, config.DbHost, config.DbPort, config.DbName)
	db, err := gorm.Open(gmsyql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Error(errors.New("数据库连接失败, " + err.Error()))
		return
	}
	sqlDB, err := db.DB()
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)
	DB = db
	DB.AutoMigrate(&model.Workflow{}, &model.User{})

}
