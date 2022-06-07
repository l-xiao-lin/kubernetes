package dao

import (
	"errors"
	"github.com/wonderivan/logger"
	"lubanKubernets/model"
	"lubanKubernets/mysql"
)

var Workflow workflow

type workflow struct {
}

type WorkflowResp struct {
	Items []model.Workflow
	Total int
}

//获取workflow列表
func (w *workflow) GetList(fileName string, limit, page int) (workflowResp *WorkflowResp, err error) {
	startSet := (page - 1) * limit

	var workflowList []model.Workflow
	tx := mysql.DB.Debug().Where("name LIKE ?", "%"+fileName+"%").Limit(limit).Offset(startSet).Order("id desc").Find(&workflowList)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error(errors.New("获取workflow列表失败, " + err.Error()))
		return nil, errors.New("获取workflow列表失败, " + err.Error())
	}
	workflowResp = &WorkflowResp{
		Items: workflowList,
		Total: len(workflowList),
	}
	return workflowResp, nil

}

//获取workflow单条信息

func (w *workflow) GetWorkflowById(id int) (workflow *model.Workflow, err error) {
	workflow = &model.Workflow{ID: id}
	tx := mysql.DB.Find(workflow)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error(errors.New("获取workflow详情失败, " + err.Error()))
		return nil, errors.New("获取workflow详情失败, " + err.Error())
	}
	return workflow, nil
}

//创建workflow

func (w *workflow) Add(workflow *model.Workflow) (err error) {
	tx := mysql.DB.Create(workflow)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error(errors.New("创建workflow失败, " + err.Error()))
		return errors.New("创建workflow失败, " + err.Error())
	}
	return nil
}

//删除workflow

func (w *workflow) DelById(id int) (err error) {
	workflow := model.Workflow{ID: id}
	tx := mysql.DB.Delete(&workflow)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error(errors.New("删除workflow失败, " + err.Error()))
		return errors.New("删除workflow失败, " + err.Error())
	}
	return nil
}
