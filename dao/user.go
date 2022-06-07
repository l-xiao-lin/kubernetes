package dao

import (
	"errors"
	"github.com/wonderivan/logger"
	"lubanKubernets/model"
	"lubanKubernets/mysql"
)

var User user

type user struct {
}

func (u *user) RegisterUser(user *model.User) (err error) {
	tx := mysql.DB.Create(user)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error(errors.New("用户注册失败, " + err.Error()))
		return errors.New("用户注册失败, " + err.Error())
	}
	return nil
}

func (u *user) Login(username, password string) (err error) {
	var user model.User
	row := mysql.DB.Where("username=? AND password=?", username, password).Find(&user).RowsAffected
	if row != 1 {
		logger.Error(errors.New("用户名或密码错误, " + err.Error()))
		return errors.New("用户名或密码错误, " + err.Error())
	}
	return nil
}
