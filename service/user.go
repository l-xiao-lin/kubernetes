package service

import (
	"lubanKubernets/dao"
	"lubanKubernets/model"
	"lubanKubernets/utils"
)

var User user

type user struct {
}

func (u *user) RegisterUser(user *model.User) (err error) {
	err = dao.User.RegisterUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (u *user) Login(username, password string) (token string, err error) {

	err = dao.User.Login(username, password)
	if err != nil {
		return "", err
	}
	tokenString, err := utils.GenToken(username)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
