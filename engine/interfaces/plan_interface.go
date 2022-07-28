package interfaces

import (
	"pk/common/conn"
	"pk/common/orm"
)

type PlanInterface interface {
	// 获取ID
	GetId() string
	// 获取排课计划名称
	GetName() string
	// 获取排课计划创建者信息
	GetPlayer() (PlayerInterface, error)
	// 获取排课计划篇章
	GetSheets(page int, limit int, wheres ...*orm.Where) (*conn.Page, error)

	// 检查排课计划是否已存在
	IsExist() (bool, error)

	// 改变排课计划名称
	ChangeName(newName string) error
	// 改变排课计划序号
	ChangeSort(newSort int) error
}