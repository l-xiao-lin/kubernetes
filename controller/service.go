package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"lubanKubernets/service"
)

var Servicev1 servicev1

type servicev1 struct {
}

func (s *servicev1) GetServiceList(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	err := c.Bind(params)
	if err != nil {
		logger.Error(errors.New("参数bind失败, " + err.Error()))
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Servicev1.GetServiceList(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取Service列表成功",
		"data": data,
	})
}

func (s *servicev1) CreateService(c *gin.Context) {
	params := new(service.ServiceCreate)
	err := c.ShouldBindJSON(params)
	if err != nil {
		logger.Error(errors.New("参数bind失败," + err.Error()))
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err = service.Servicev1.CreateService(params)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "创建Service成功",
		"data": nil,
	})
}

func (s *servicev1) GetServiceDetail(c *gin.Context) {
	params := new(struct {
		ServiceName string `form:"service_name"`
		Namespace   string `form:"namespace"`
	})
	err := c.Bind(params)
	if err != nil {
		logger.Error(errors.New("参数bind失败, " + err.Error()))
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Servicev1.GetServiceDetail(params.ServiceName, params.Namespace)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	jsonStr, _ := json.Marshal(data)

	c.JSON(200, gin.H{
		"msg":  "获取Service详情成功",
		"data": string(jsonStr),
	})

}

func (s *servicev1) DeleteService(c *gin.Context) {
	params := new(struct {
		ServiceName string `json:"service_name"`
		Namespace   string `json:"namespace"`
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
	err = service.Servicev1.DeleteService(params.ServiceName, params.Namespace)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "删除Service成功",
		"data": nil,
	})

}

func (s *servicev1) UpdateService(c *gin.Context) {
	params := new(struct {
		Content   string `json:"content"`
		Namespace string `json:"namespace"`
	})
	err := c.ShouldBindJSON(params)
	if err != nil {
		logger.Error(errors.New("参数bind错误, " + err.Error()))
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err = service.Servicev1.UpdateService(params.Content, params.Namespace)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "更新Service成功",
		"data": nil,
	})

}
