package interfaces

import (
	"pk/common/conn"
	"pk/common/orm"
)

// 被称作玩家需要被实现的功能
type PlayerInterface interface {
	// 获取ID
	GetId() string
	// 获取玩家名称
	GetName() string
	// 获取单位信息
	GetUnit() (UnitInterface, error)
	// 获取排课计划列表
	GetPlans(page, limit int, wheres ...*orm.Where) (*conn.Page, error)

	// 改变玩家登录密码
	ChangePassword(newPassword string) error
	// 改变玩家名称
	ChangeName(newName string) error
	// 改变玩家序号
	ChangeSort(newSort int) error

	// 创建排课计划
	CreatePlan(planName string, sort ...int) (PlanInterface, error)
}