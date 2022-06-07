package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var Deployment deployment

type deployment struct {
}

type DeploymentCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Label         map[string]string `json:"label"`
	Replicas      int32             `json:"replicas"`
	Image         string            `json:"image"`
	ContainerPort int32             `json:"container_port"`
	HealthCheck   bool              `json:"health_check"`
	HealthPath    string            `json:"health_path"`
	Cpu           string            `json:"cpu"`
	Memory        string            `json:"memory"`
}

type DeploymentsResp struct {
	Items []appsv1.Deployment `json:"items"`
	Total int                 `json:"total"`
}

type DeploymentsNp struct {
	Namespace string
	Total     int
}

func (d *deployment) toCell(std []appsv1.Deployment) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = deploymentCell(std[i])
	}
	return cells

}

func (d *deployment) fromCell(cells []DataCell) []appsv1.Deployment {
	deployments := make([]appsv1.Deployment, len(cells))
	for i := range cells {
		deployments[i] = appsv1.Deployment(cells[i].(deploymentCell))
	}
	return deployments
}

//创建Deployment

func (d *deployment) CreateDeployment(data *DeploymentCreate) (err error) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
			Labels:    data.Label,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &data.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: data.Label,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: data.Label,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  data.Name,
							Image: data.Image,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: data.ContainerPort,
								},
							},
						},
					},
				},
			},
		},
	}

	if data.HealthCheck {
		deployment.Spec.Template.Spec.Containers[0].ReadinessProbe = &corev1.Probe{
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
		deployment.Spec.Template.Spec.Containers[0].LivenessProbe = &corev1.Probe{
			Handler: corev1.Handler{
				HTTPGet: &corev1.HTTPGetAction{
					Path: data.HealthPath,
					Port: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
			InitialDelaySeconds: 5,
			TimeoutSeconds:      5,
			PeriodSeconds:       5,
		}
	}
	deployment.Spec.Template.Spec.Containers[0].Resources.Limits = corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}

	deployment.Spec.Template.Spec.Containers[0].Resources.Requests = corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse(data.Cpu),
		corev1.ResourceMemory: resource.MustParse(data.Memory),
	}

	_, err = K8s.ClientSet.AppsV1().Deployments(data.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建Deployment失败, " + err.Error()))
		return errors.New("创建Deployment失败, " + err.Error())
	}
	return nil
}

//删除Deployment

func (d *deployment) DeleteDeployment(deploymentName, namespace string) (err error) {
	err = K8s.ClientSet.AppsV1().Deployments(namespace).Delete(context.TODO(), deploymentName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Deployment失败, " + err.Error()))
		return errors.New("删除Deployment失败, " + err.Error())
	}
	return nil
}

//获取Deployment列表

func (d *deployment) GetDeploymentList(filterName, namespace string, limit, page int) (deploymentsResp *DeploymentsResp, err error) {

	deploymentList, err := K8s.ClientSet.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Deployment列表失败, " + err.Error()))
		return nil, errors.New("获取Deployment列表失败, " + err.Error())
	}
	selectableData := dataSelector{
		GenericDataList: d.toCell(deploymentList.Items),
		dataSelectQuery: &DataSelectQuery{
			FilterQuery: &FilterQuery{
				Name: filterName,
			},
			PaginateQuery: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}

	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	filtered.Sort().Paginate()

	deployments := d.fromCell(filtered.GenericDataList)

	deploymentsResp = &DeploymentsResp{
		Items: deployments,
		Total: total,
	}
	return deploymentsResp, nil

}

//获取Deployment详情

func (d *deployment) GetDeploymentDetail(deploymentName, namespace string) (deployment *appsv1.Deployment, err error) {
	deployment, err = K8s.ClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Deployment详情失败, " + err.Error()))
		return nil, errors.New("获取Deployment详情失败, " + err.Error())
	}

	return deployment, nil
}

//修改Deployment

func (d *deployment) UpdateDeployment(content, namespace string) (err error) {
	var deployment appsv1.Deployment
	err = json.Unmarshal([]byte(content), &deployment)
	if err != nil {
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}

	_, err = K8s.ClientSet.AppsV1().Deployments(namespace).Update(context.TODO(), &deployment, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Deployment失败, " + err.Error()))
		return errors.New("更新Deployment失败, " + err.Error())
	}
	return nil

}

//设置Scale副本数

func (d *deployment) ScaleDeployment(deploymentName, namespace string, scaleNum int32) (replicas int32, err error) {
	scale, err := K8s.ClientSet.AppsV1().Deployments(namespace).GetScale(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Scale失败, " + err.Error()))
		return 0, errors.New("获取Scale失败, " + err.Error())
	}
	scale.Spec.Replicas = scaleNum

	newScale, err := K8s.ClientSet.AppsV1().Deployments(namespace).UpdateScale(context.TODO(), deploymentName, scale, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Scale失败, " + err.Error()))
		return 0, errors.New("更新Scale失败, " + err.Error())
	}
	return newScale.Spec.Replicas, nil
}

//获取每个namespace中的Deployment数量

func (d *deployment) GetDeploymentsPerNp() (deploymentsNps []*DeploymentsNp, err error) {
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})

	for _, namespace := range namespaceList.Items {
		deploymentList, err := K8s.ClientSet.AppsV1().Deployments(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logger.Error(errors.New("获取Namespace失败, " + err.Error()))
			return nil, errors.New("获取Namespace失败, " + err.Error())
		}
		deploymentsNp := &DeploymentsNp{
			Namespace: namespace.Name,
			Total:     len(deploymentList.Items),
		}
		deploymentsNps = append(deploymentsNps, deploymentsNp)
	}
	return deploymentsNps, nil
}
