package model

type User struct {
	Id       int    `gorm:"primaryKey"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (User) TableName() string {
	return "user"
}
