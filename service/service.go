package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var Servicev1 servicev1

type servicev1 struct {
}

type ServicesResp struct {
	Items []corev1.Service
	Total int
}

type ServiceCreate struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Label         map[string]string `json:"label"`
	Port          int32             `json:"port"`
	NodePort      int32             `json:"node_port"`
	ContainerPort int32             `json:"container_port"`
	Type          string            `json:"type"`
}

func (s *servicev1) ToCells(std []corev1.Service) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = serviceCell(std[i])
	}
	return cells
}

func (s *servicev1) FromCells(cells []DataCell) []corev1.Service {
	services := make([]corev1.Service, len(cells))
	for i := range cells {
		services[i] = corev1.Service(cells[i].(serviceCell))
	}
	return services
}

//获取Services列表

func (s *servicev1) GetServiceList(filterName, namespace string, limit, page int) (servicesResp *ServicesResp, err error) {
	serviceList, err := K8s.ClientSet.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Services列表失败, " + err.Error()))
		return nil, errors.New("获取Services列表失败, " + err.Error())
	}
	selectableData := dataSelector{
		GenericDataList: s.ToCells(serviceList.Items),
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

	services := s.FromCells(filtered.GenericDataList)

	servicesResp = &ServicesResp{
		Items: services,
		Total: total,
	}
	return servicesResp, nil
}

//创建Service

func (s *servicev1) CreateService(data *ServiceCreate) (err error) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: data.Label,
			Ports: []corev1.ServicePort{
				{
					Protocol: corev1.ProtocolTCP,
					Port:     data.Port,
					TargetPort: intstr.IntOrString{
						Type:   0,
						IntVal: data.ContainerPort,
					},
				},
			},
		},
	}

	_, err = K8s.ClientSet.CoreV1().Services(data.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建Service失败, " + err.Error()))
		return errors.New("创建Service失败, " + err.Error())
	}
	return nil
}

//获取Service详情

func (s *servicev1) GetServiceDetail(serviceName, namespace string) (service *corev1.Service, err error) {
	service, err = K8s.ClientSet.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Service失败, " + err.Error()))
		return nil, errors.New("获取Service失败, " + err.Error())
	}
	return service, nil
}

//删除Service详情

func (s *servicev1) DeleteService(serviceName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Services(namespace).Delete(context.TODO(), serviceName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Service失败, " + err.Error()))
		return errors.New("删除Service失败, " + err.Error())
	}
	return nil
}

//更新Service

func (s *servicev1) UpdateService(content, namespace string) (err error) {
	service := &corev1.Service{}

	err = json.Unmarshal([]byte(content), service)
	if err != nil {
		logger.Error(errors.New("Content内容反序列失败, " + err.Error()))
		return errors.New("Content内容反序列失败, " + err.Error())
	}
	_, err = K8s.ClientSet.CoreV1().Services(namespace).Update(context.TODO(), service, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Service失败, " + err.Error()))
		return errors.New("更新Service失败, " + err.Error())
	}
	return nil
}
