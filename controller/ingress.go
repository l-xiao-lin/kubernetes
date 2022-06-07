package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"lubanKubernets/service"
)

var Ingress ingress

type ingress struct {
}

func (i *ingress) GetIngressList(c *gin.Context) {
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
	data, err := service.Ingress.GetIngressList(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取Ingress列表成功",
		"data": data,
	})
}

func (i *ingress) CreateIngress(c *gin.Context) {
	params := new(service.IngressCreate)
	err := c.ShouldBindJSON(params)
	if err != nil {
		logger.Error(errors.New("参数bind失败, " + err.Error()))
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err = service.Ingress.CreateIngress(params)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "创建Ingress成功",
		"data": nil,
	})

}

func (i *ingress) GetIngressDetail(c *gin.Context) {
	params := new(struct {
		IngressName string `form:"ingress_name"`
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
	data, err := service.Ingress.GetIngressDetail(params.IngressName, params.Namespace)

	jsonStr, _ := json.Marshal(data)

	c.JSON(200, gin.H{
		"msg":  "获取Ingress详情成功",
		"data": string(jsonStr),
	})

}

func (i *ingress) UpdateIngress(c *gin.Context) {
	params := new(struct {
		Content   string `json:"content"`
		Namespace string `json:"namespace"`
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
	err = service.Ingress.UpdateIngress(params.Content, params.Namespace)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "Ingress已更新",
		"data": nil,
	})
}

func (i *ingress) DeleteIngress(c *gin.Context) {
	params := new(struct {
		IngressName string `json:"ingress_name"`
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
	err = service.Ingress.DeleteIngress(params.IngressName, params.Namespace)
	c.JSON(200, gin.H{
		"msg":  "删除Ingress成功",
		"data": nil,
	})

}
