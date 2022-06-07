package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"lubanKubernets/service"
)

var Workflow workflow

type workflow struct {
}

func (w *workflow) GetList(c *gin.Context) {
	params := new(struct {
		FilerName string `form:"filter_name"`
		Limit     int    `form:"limit"`
		Page      int    `form:"page"`
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
	data, err := service.Workflow.GetList(params.FilerName, params.Limit, params.Page)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取workflow列表成功",
		"data": data,
	})

}

func (w *workflow) GetWorkflowById(c *gin.Context) {
	params := new(struct {
		Id int `form:"id"`
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
	data, err := service.Workflow.GetWorkflowById(params.Id)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取workflow详情成功",
		"data": data,
	})

}

func (w *workflow) CreateWorkflow(c *gin.Context) {
	params := new(service.WorkflowCreate)
	err := c.ShouldBindJSON(params)
	if err != nil {
		logger.Error(errors.New("参数bind失败, " + err.Error()))
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err = service.Workflow.CreateWorkflow(params)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "创建workflow成功",
		"data": nil,
	})
}

func (w *workflow) DeleteWorkflow(c *gin.Context) {
	params := new(struct {
		Id int `json:"id"`
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
	err = service.Workflow.DeleteWorkflow(params.Id)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "删除workflow成功",
		"data": nil,
	})

}
