package main

import (
	"github.com/gin-gonic/gin"
	"lubanKubernets/config"
	"lubanKubernets/controller"
	"lubanKubernets/middleware"
	"lubanKubernets/mysql"
	"lubanKubernets/service"
)

func main() {
	r := gin.Default()

	//初始化数据库
	mysql.Init()

	//jwt token验证
	r.Use(middleware.JwtAuth())

	//跨域配置
	r.Use(middleware.Cors())

	//加载路由
	controller.Router.InitApiRouter(r)

	//初始化k8s client
	service.K8s.Init()

	r.Run(config.ListenAddr)
}
