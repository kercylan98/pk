package interfaces

import (
	"pk/common/conn"
	"pk/common/orm"
)

type UnitGradeInterface interface {
	// 获取ID
	GetId() string
	// 获取年级名称
	GetName() string

	// 获取单位信息
	GetUnit() (UnitInterface, error)
	// 获取年级下的班级
	GetClasses(page, limit int, wheres ...*orm.Where) (*conn.Page, error)

	// 改变年级名称
	ChangeName(newName string) error
	// 改变年级序号
	ChangeSort(newSort int) error
}
