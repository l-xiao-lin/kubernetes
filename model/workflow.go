package model

type Workflow struct {
	ID         int `gorm:"primaryKey"`
	Name       string
	Namespace  string
	Replicas   int32
	Deployment string
	Service    string
	Ingress    string
	Type       string
}

func (Workflow) TableName() string {
	return "workflow"
}
