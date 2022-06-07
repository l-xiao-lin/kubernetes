package service

import (
	"lubanKubernets/dao"
	"lubanKubernets/model"
)

var Workflow workflow

type workflow struct {
}

type WorkflowCreate struct {
	Name          string                `json:"name"`
	Namespace     string                `json:"namespace"`
	Label         map[string]string     `json:"label"`
	Replicas      int32                 `json:"replicas"`
	Image         string                `json:"image"`
	ContainerPort int32                 `json:"container_port"`
	HealthCheck   bool                  `json:"health_check"`
	HealthPath    string                `json:"health_path"`
	Cpu           string                `json:"cpu"`
	Memory        string                `json:"memory"`
	Port          int32                 `json:"port"`
	NodePort      int32                 `json:"node_port"`
	Type          string                `json:"type"`
	Hosts         map[string][]httpPath `json:"hosts"`
}

//查询workflow分页列表
func (w *workflow) GetList(filerName string, limit, page int) (data *dao.WorkflowResp, err error) {
	data, err = dao.Workflow.GetList(filerName, limit, page)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//查询workflow详情

func (w *workflow) GetWorkflowById(id int) (data *model.Workflow, err error) {
	data, err = dao.Workflow.GetWorkflowById(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//创建workflow

func (w *workflow) CreateWorkflow(data *WorkflowCreate) (err error) {
	//添加k8s资源

	err = CreateWorkflowRes(data)
	if err != nil {
		return err
	}

	//添加mysql数据
	var ingressName string
	if data.Type == "Ingress" {
		ingressName = GetIngressName(data.Name)
	} else {
		ingressName = ""
	}

	workflow := model.Workflow{
		Name:       data.Name,
		Namespace:  data.Namespace,
		Replicas:   data.Replicas,
		Deployment: data.Name,
		Service:    GetServiceName(data.Name),
		Ingress:    ingressName,
		Type:       data.Type,
	}
	err = dao.Workflow.Add(&workflow)
	if err != nil {
		return err
	}

	return nil
}

func CreateWorkflowRes(data *WorkflowCreate) (err error) {
	//创建deployment资源
	deployment := DeploymentCreate{
		Name:          data.Name,
		Namespace:     data.Namespace,
		Label:         data.Label,
		Replicas:      data.Replicas,
		Image:         data.Image,
		ContainerPort: data.ContainerPort,
		HealthCheck:   data.HealthCheck,
		HealthPath:    data.HealthPath,
		Cpu:           data.Cpu,
		Memory:        data.Memory,
	}
	err = Deployment.CreateDeployment(&deployment)
	if err != nil {
		return err
	}
	//创建svc资源
	var ServiceType string
	if data.Type == "NodePort" {
		ServiceType = data.Type
	} else {
		ServiceType = "ClusterIP"
	}

	service := ServiceCreate{
		Name:          GetServiceName(data.Name),
		Namespace:     data.Namespace,
		Label:         data.Label,
		Port:          data.Port,
		NodePort:      data.NodePort,
		ContainerPort: data.ContainerPort,
		Type:          ServiceType,
	}
	err = Servicev1.CreateService(&service)
	if err != nil {
		return err
	}

	//创建ingress资源
	if data.Type == "Ingress" {
		ingress := IngressCreate{
			Name:      GetIngressName(data.Name),
			Namespace: data.Namespace,
			Label:     data.Label,
			Hosts:     data.Hosts,
		}
		err = Ingress.CreateIngress(&ingress)
		if err != nil {
			return err
		}
		return nil
	}
	return nil

}

func GetServiceName(name string) string {
	return name + "-svc"
}

func GetIngressName(name string) string {
	return name + "-ing"
}

//删除workflow

func (w *workflow) DeleteWorkflow(id int) (err error) {

	//先获取数据库条目
	workflow, err := dao.Workflow.GetWorkflowById(id)

	//删除k8s资源
	err = DeleteWorkflowRes(workflow)
	if err != nil {
		return err
	}

	//删除数据库数据
	err = dao.Workflow.DelById(id)
	if err != nil {
		return err
	}
	return nil

}

func DeleteWorkflowRes(workflow *model.Workflow) (err error) {
	//删除deployment资源
	err = Deployment.DeleteDeployment(workflow.Deployment, workflow.Namespace)
	if err != nil {
		return err
	}

	//删除service资源

	err = Servicev1.DeleteService(workflow.Service, workflow.Namespace)
	if err != nil {
		return err
	}

	//删除ingress资源
	if workflow.Type == "Ingress" {
		err = Ingress.DeleteIngress(workflow.Ingress, workflow.Namespace)
		if err != nil {
			return err
		}
	}

	return nil

}
