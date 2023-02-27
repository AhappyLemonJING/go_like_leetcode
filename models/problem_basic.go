package models

import (
	"gorm.io/gorm"
)

type ProblemBasic struct {
	gorm.Model
	Identity          string             `gorm:"column:identity;type:varchar(36);" json:"identity"` // 问题的唯一标识
	ProblemCategories []*ProblemCategory `gorm:"foreignKey:problem_id;references:id"`
	Title             string             `gorm:"column:title;type:varchar(255);" json:"title"` // 题目的标题
	Content           string             `gorm:"column:content;type:text;" json:"content"`     // 题目正文描述
	MaxMem            int                `gorm:"column:max_mem;type:int;" json:"max_mem"`
	MaxRuntime        int                `gorm:"column:max_runtime;type:int;" json:"max_runtime"`
	TestCase          []*TestCase        `gorm:"foreignKey:problem_identity;references:identity"`
	PassNum           int64              `gorm:"column:pass_num;type:int(11);" json:"pass_num"`     // 通过个数
	SubmitNum         int64              `gorm:"column:submit_num;type:int(11);" json:"submit_num"` // 提交次数
}

func (table *ProblemBasic) TableName() string {
	return "problem_basic"
}

func GetProblemList(keyword string, categoryIdentity string) *gorm.DB {
	tx := DB.Model(new(ProblemBasic)).Preload("ProblemCategories").Preload("ProblemCategories.CategoryBasic").Where("title like ? OR content like ?", "%"+keyword+"%", "%"+keyword+"%")

	if categoryIdentity != "" {
		tx.Joins("RIGHT JOIN problem_category pc on pc.problem_id = problem_basic.id").Where("pc.category_id = (SELECT cb.id FROM category_basic cb WHERE cb.identity = ?)", categoryIdentity)
	}
	return tx
}
