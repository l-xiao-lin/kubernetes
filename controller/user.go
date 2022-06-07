package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"lubanKubernets/model"
	"lubanKubernets/service"
)

var User user

type user struct {
}

func (u *user) RegisterUser(c *gin.Context) {
	params := new(model.User)
	err := c.ShouldBindJSON(params)
	if err != nil {
		logger.Error(errors.New("参数bind失败, " + err.Error()))
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	service.User.RegisterUser(params)
	c.JSON(200, gin.H{
		"msg":  "注册成功",
		"data": nil,
	})
}

func (u *user) Login(c *gin.Context) {
	params := new(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	})
	err := c.ShouldBindJSON(params)
	if err != nil {
		logger.Error(errors.New("参数bind失败, " + err.Error()))
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.User.Login(params.Username, params.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "用户成功",
		"data": data,
	})

}
