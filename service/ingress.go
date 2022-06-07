package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Ingress ingress

type ingress struct {
}

type IngressResp struct {
	Items []v1.Ingress
	Total int
}

type IngressCreate struct {
	Name      string                `json:"name"`
	Namespace string                `json:"namespace"`
	Label     map[string]string     `json:"label"`
	Hosts     map[string][]httpPath `json:"hosts"`
}

type httpPath struct {
	Path        string      `json:"path"`
	PathType    v1.PathType `json:"path_type"`
	ServiceName string      `json:"service_name"`
	ServicePort int32       `json:"service_port"`
}

func (i *ingress) toCells(std []v1.Ingress) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = ingressCell(std[i])
	}
	return cells
}

func (i *ingress) fromCells(cells []DataCell) []v1.Ingress {
	ingresses := make([]v1.Ingress, len(cells))
	for i := range cells {
		ingresses[i] = v1.Ingress(cells[i].(ingressCell))
	}
	return ingresses
}

//获取ingress列表

func (i *ingress) GetIngressList(filterName, namespace string, limit, page int) (ingressResp *IngressResp, err error) {
	ingressList, err := K8s.ClientSet.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取ingressList失败, " + err.Error()))
		return nil, errors.New("获取ingressList失败, " + err.Error())
	}
	selectableData := &dataSelector{
		GenericDataList: i.toCells(ingressList.Items),
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
	data := filtered.Sort().Paginate()

	ingresses := i.fromCells(data.GenericDataList)

	ingressResp = &IngressResp{
		Items: ingresses,
		Total: total,
	}
	return ingressResp, nil

}

//创建Ingress

func (i *ingress) CreateIngress(data *IngressCreate) (err error) {
	var httpIngressPaths []v1.HTTPIngressPath
	var ingressRules []v1.IngressRule

	ingress := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      data.Name,
			Namespace: data.Namespace,
		},
		Status: v1.IngressStatus{},
	}

	for key, value := range data.Hosts {
		ir := v1.IngressRule{
			Host: key,
			IngressRuleValue: v1.IngressRuleValue{
				HTTP: &v1.HTTPIngressRuleValue{
					Paths: nil,
				},
			},
		}

		for _, httpPath := range value {
			hip := v1.HTTPIngressPath{
				Path:     httpPath.Path,
				PathType: &httpPath.PathType,
				Backend: v1.IngressBackend{
					Service: &v1.IngressServiceBackend{
						Name: httpPath.ServiceName,
						Port: v1.ServiceBackendPort{
							Number: httpPath.ServicePort,
						},
					},
				},
			}
			httpIngressPaths = append(httpIngressPaths, hip)
		}
		ir.IngressRuleValue.HTTP.Paths = httpIngressPaths

		ingressRules = append(ingressRules, ir)
	}
	ingress.Spec.Rules = ingressRules

	_, err = K8s.ClientSet.NetworkingV1().Ingresses(data.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		logger.Error(errors.New("创建ingress失败, " + err.Error()))
		return errors.New("创建ingress失败, " + err.Error())
	}
	return nil

}

func (i *ingress) GetIngressDetail(ingressName, namespace string) (ingress *v1.Ingress, err error) {

	ingress, err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingressName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Ingress详情失败, " + err.Error()))
		return nil, errors.New("获取Ingress详情失败, " + err.Error())
	}
	return ingress, nil

}

func (i *ingress) UpdateIngress(content, namespace string) (err error) {
	ingress := &v1.Ingress{}
	err = json.Unmarshal([]byte(content), &ingress)
	if err != nil {
		logger.Error(errors.New("content反序列化失败," + err.Error()))
		return errors.New("content反序列化失败," + err.Error())
	}
	_, err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Update(context.TODO(), ingress, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Ingress失败, " + err.Error()))
		return errors.New("更新Ingress失败, " + err.Error())
	}
	return nil

}

func (i *ingress) DeleteIngress(ingressName, namespace string) (err error) {
	err = K8s.ClientSet.NetworkingV1().Ingresses(namespace).Delete(context.TODO(), ingressName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Ingress失败, " + err.Error()))
		return errors.New("删除Ingress失败, " + err.Error())
	}
	return nil

}
