package models

import "gorm.io/gorm"

type UserBasic struct {
	gorm.Model
	Identity  string `gorm:"column:identity;type:varchar(36);" json:"identity"`       // 用户的唯一标识
	Name      string `gorm:"column:name;type:varchar(100);" json:"name"`              // 姓名
	Password  string `gorm:"column:password;type:varchar(32);" json:"password"`       // 密码
	Phone     string `gorm:"column:phone;type:varchar(20);" json:"phone"`             // 电话
	Mail      string `gorm:"column:mail;type:varchar(100);" json:"mail"`              // 邮箱
	PassNum   int64  `gorm:"column:finish_problem_num;type:int(11);" json:"pass_num"` // 通过个数
	SubmitNum int64  `gorm:"column:submit_num;type:int(11);" json:"submit_num"`       // 提交次数
	IsAdmin   int    `gorm:"column:is_admin;type:tinyint(1);" json:"is_admin"`
}

func (table *UserBasic) TableName() string {
	return "user_basic"
}
