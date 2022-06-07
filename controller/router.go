package controller

import "github.com/gin-gonic/gin"

var Router router

type router struct {
}

func (r *router) InitApiRouter(router *gin.Engine) {
	router.GET("/demo", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg":  "success",
			"data": nil,
		})
	})
	//pod操作
	router.GET("/api/k8s/pods", Pod.GetPodList)
	router.GET("/api/k8s/pod/detail", Pod.GetPodDetail)
	router.DELETE("/api/k8s/pod/del", Pod.DeletePod)
	router.PUT("/api/k8s/pod/update", Pod.UpdatePod)
	router.GET("/api/k8s/pod/container", Pod.GetContainers)
	router.GET("/api/k8s/pod/log", Pod.GetContainerLog)
	router.GET("/api/k8s/pod/numnp", Pod.GetPodsPerNp)

	//deployment操作
	router.GET("/api/k8s/deployments", Deployment.GetDeploymentList)
	router.GET("/api/k8s/deployment/detail", Deployment.GetDeploymentDetail)
	router.POST("/api/k8s/deployment/create", Deployment.CreateDeployment)
	router.DELETE("/api/k8s/deployment/del", Deployment.DeleteDeployment)
	router.PUT("/api/k8s/deployment/scale", Deployment.ScaleDeployment)
	router.GET("/api/k8s/deployment/numnp", Deployment.GetDeploymentsPerNp)
	router.PUT("/api/k8s/deployment/update", Deployment.UpdateDeployment)

	//service操作
	router.GET("/api/k8s/services/", Servicev1.GetServiceList)
	router.POST("/api/k8s/service/create", Servicev1.CreateService)
	router.GET("/api/k8s/service/detail", Servicev1.GetServiceDetail)
	router.DELETE("/api/k8s/service/del", Servicev1.DeleteService)
	router.PUT("/api/k8s/service/update", Servicev1.UpdateService)

	//ingress操作
	router.GET("/api/k8s/ingresses", Ingress.GetIngressList)
	router.GET("/api/k8s/ingress/detail", Ingress.GetIngressDetail)
	router.POST("/api/k8s/ingress/create", Ingress.CreateIngress)
	router.PUT("/api/k8s/ingress/update", Ingress.UpdateIngress)
	router.DELETE("/api/k8s/ingress/del", Ingress.DeleteIngress)

	//workflow操作
	router.GET("/api/k8s/workflows", Workflow.GetList)
	router.GET("/api/k8s/workflow/detail", Workflow.GetWorkflowById)
	router.POST("/api/k8s/workflow/create", Workflow.CreateWorkflow)
	router.DELETE("/api/k8s/workflow/del", Workflow.DeleteWorkflow)

	//注册
	router.POST("/api/k8s/register", User.RegisterUser)
	router.POST("/api/k8s/login", User.Login)

}
