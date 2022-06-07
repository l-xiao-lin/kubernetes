package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/wonderivan/logger"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"lubanKubernets/config"
)

var Pod pod

type pod struct {
}

type PodsResp struct {
	Items []corev1.Pod
	Total int
}

type PodsPerNp struct {
	Namespace string
	Total     int
}

func (p *pod) toCells(std []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = podCell(std[i])
	}
	return cells
}

func (p *pod) fromCell(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}

//获取Pod列表信息

func (p *pod) GetPodList(filterName, namespace string, limit, page int) (podsResp *PodsResp, err error) {
	podList, err := K8s.ClientSet.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Pod列表失败, " + err.Error()))
		return nil, errors.New("获取Pod列表失败, " + err.Error())
	}
	selectableData := dataSelector{
		GenericDataList: p.toCells(podList.Items),
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
	pods := p.fromCell(data.GenericDataList)
	podsResp = &PodsResp{
		Items: pods,
		Total: total,
	}
	return podsResp, nil

}

//获取Pod详情

func (p *pod) GetPodDetail(podName, namespace string) (pod *corev1.Pod, err error) {
	pod, err = K8s.ClientSet.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取pod详情失败, " + err.Error()))
		return nil, errors.New("获取pod详情失败, " + err.Error())
	}
	return pod, nil
}

//删除Pod

func (p *pod) DeletePod(podName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Pod失败, " + err.Error()))
		return errors.New("删除Pod失败, " + err.Error())
	}
	return nil
}

//更新Pod

func (p *pod) UpdatePod(content, namespace string) (err error) {
	pod := new(corev1.Pod)
	err = json.Unmarshal([]byte(content), &pod)
	if err != nil {
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}

	_, err = K8s.ClientSet.CoreV1().Pods(namespace).Update(context.TODO(), pod, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("Pod更新失败, " + err.Error()))
		return errors.New("Pod更新失败, " + err.Error())
	}
	return nil

}

//获取每个namespace中pod数量

func (p *pod) GetPodsPerNp() (podsPerNps []*PodsPerNp, err error) {
	namespaceList, err := K8s.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取namespace列表失败, " + err.Error()))
		return nil, errors.New("获取namespace列表失败, " + err.Error())
	}
	for _, namespace := range namespaceList.Items {
		podList, err := K8s.ClientSet.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			logger.Error(errors.New("获取pod列表失败, " + err.Error()))
			return nil, errors.New("获取pod列表失败, " + err.Error())

		}
		podsPerNp := &PodsPerNp{
			Namespace: namespace.Name,
			Total:     len(podList.Items),
		}
		podsPerNps = append(podsPerNps, podsPerNp)
	}
	return podsPerNps, nil

}

//获取Pod中Container

func (p *pod) GetContainers(podName, namespace string) (containers []string, err error) {
	pod, err := p.GetPodDetail(podName, namespace)
	if err != nil {
		return nil, err
	}

	for _, container := range pod.Spec.Containers {
		containers = append(containers, container.Name)
	}
	return containers, nil

}

//获取Pod中的日志

func (p *pod) GetContainerLog(podName, containerName, namespace string) (logs string, err error) {
	tailLine := int64(config.PodLogTailLine)
	option := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &tailLine,
	}

	rep := K8s.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, option)

	podLogs, err := rep.Stream(context.TODO())
	if err != nil {
		logger.Error(errors.New("获取podLogs失败, " + err.Error()))
		return "", errors.New("获取podLogs失败, " + err.Error())
	}

	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		logger.Error(errors.New("复制podLogs失败, " + err.Error()))
		return "", errors.New("复制podLogs失败, " + err.Error())
	}
	return buf.String(), nil

}
