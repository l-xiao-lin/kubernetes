package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
	"lubanKubernets/service"
)

var Deployment deployment

type deployment struct {
}

func (d *deployment) CreateDeployment(c *gin.Context) {
	params := new(service.DeploymentCreate)
	err := c.ShouldBindJSON(params)
	if err != nil {
		logger.Error("参数bind失败, " + err.Error())
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err = service.Deployment.CreateDeployment(params)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "创建Deployment成功",
		"data": nil,
	})

}

func (d *deployment) DeleteDeployment(c *gin.Context) {
	params := new(struct {
		DeploymentName string `json:"deployment_name"`
		Namespace      string `json:"namespace"`
	})
	err := c.ShouldBindJSON(params)
	if err != nil {
		logger.Error("参数bind失败, " + err.Error())
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err = service.Deployment.DeleteDeployment(params.DeploymentName, params.Namespace)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "成功删除Deployment",
		"data": nil,
	})
}

func (d *deployment) GetDeploymentList(c *gin.Context) {
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	err := c.Bind(params)
	if err != nil {
		logger.Error("参数bind失败, " + err.Error())
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Deployment.GetDeploymentList(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取Deployment列表成功",
		"data": data,
	})
}

func (d *deployment) GetDeploymentDetail(c *gin.Context) {
	params := new(struct {
		DeploymentName string `form:"deployment_name"`
		Namespace      string `form:"namespace"`
	})
	err := c.Bind(params)
	if err != nil {
		logger.Error("参数bind失败, " + err.Error())
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Deployment.GetDeploymentDetail(params.DeploymentName, params.Namespace)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	jsonStr, _ := json.Marshal(data)

	c.JSON(200, gin.H{
		"msg":  "获取Deployment详情成功",
		"data": string(jsonStr),
	})

}

func (d *deployment) UpdateDeployment(c *gin.Context) {
	params := new(struct {
		Content   string `json:"content"`
		Namespace string `json:"namespace"`
	})
	err := c.ShouldBindJSON(params)
	if err != nil {
		logger.Error("参数bind失败, " + err.Error())
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	err = service.Deployment.UpdateDeployment(params.Content, params.Namespace)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "更新Deployment成功",
		"data": nil,
	})
}

func (d *deployment) ScaleDeployment(c *gin.Context) {
	params := new(struct {
		DeploymentName string `json:"deployment_name"`
		Namespace      string `json:"namespace"`
		ScaleNum       int32  `json:"scale_num"`
	})
	err := c.ShouldBindJSON(params)
	if err != nil {
		logger.Error("参数bind失败, " + err.Error())
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Deployment.ScaleDeployment(params.DeploymentName, params.Namespace, params.ScaleNum)
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "Scale副本数更新成功",
		"data": data,
	})

}

func (d *deployment) GetDeploymentsPerNp(c *gin.Context) {

	data, err := service.Deployment.GetDeploymentsPerNp()
	if err != nil {
		c.JSON(500, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":  "获取每个namespace中的deployment数量",
		"data": data,
	})

}
