package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"lubanKubernets/service"
)

var Pod pod

type pod struct {
}

func (p *pod) GetPodList(c *gin.Context) {
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
	data, err := service.Pod.GetPodList(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取pod列表成功",
		"data": data,
	})

}

func (p *pod) GetPodDetail(c *gin.Context) {
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
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
	data, err := service.Pod.GetPodDetail(params.PodName, params.Namespace)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	jsonStr, _ := json.Marshal(data)

	c.JSON(200, gin.H{
		"msg":  "获取pod详情成功",
		"data": string(jsonStr),
	})

}

func (p *pod) DeletePod(c *gin.Context) {
	params := new(struct {
		PodName   string `json:"pod_name"`
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
	err = service.Pod.DeletePod(params.PodName, params.Namespace)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"msg":  "成功删除pod",
		"data": nil,
	})

}

func (p *pod) UpdatePod(c *gin.Context) {
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
	err = service.Pod.UpdatePod(params.Content, params.Namespace)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	c.JSON(200, gin.H{
		"msg":  "更新pod信息成功",
		"data": nil,
	})

}

func (p *pod) GetPodsPerNp(c *gin.Context) {
	data, err := service.Pod.GetPodsPerNp()
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "成功获取每个namespace中的pod数量",
		"data": data,
	})

}

func (p *pod) GetContainers(c *gin.Context) {
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
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
	data, err := service.Pod.GetContainers(params.PodName, params.Namespace)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取container成功",
		"data": data,
	})

}

func (p *pod) GetContainerLog(c *gin.Context) {
	params := new(struct {
		PodName       string `form:"pod_name"`
		ContainerName string `form:"container_name"`
		Namespace     string `form:"namespace"`
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
	data, err := service.Pod.GetContainerLog(params.PodName, params.ContainerName, params.Namespace)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "成功获取podLog日志",
		"data": data,
	})
}
