package interfaces

import (
	"pk/common/conn"
	"pk/common/orm"
)

type UnitInterface interface {
	// 获取单位ID
	GetId() string
	// 获取单位名称
	GetName() string
	// 获取单位下的玩家
	GetPlayers(page, limit int, wheres ...*orm.Where) (*conn.Page, error)
	// 获取单位下的年级
	GetGrades(page, limit int, wheres ...*orm.Where) (*conn.Page, error)
	// 获取单位下的教学场地
	GetSites(page, limit int, wheres ...*orm.Where) (*conn.Page, error)
	// 获取单位下的教师
	GetTeachers(page, limit int, wheres ...*orm.Where) (*conn.Page, error)

	// 增加年级
	AddGrade(gradeName string, sort ...int) (UnitGradeInterface, error)
	// 增加教学场地
	AddSite(siteName string, sort ...int) (UnitSiteInterface, error)
	// 增加教师
	AddTeacher(teacherName string, sort ...int) (UnitTeacherInterface, error)

	// 改变单位名称
	ChangeName(newName string) error
	// 改变单位序号
	ChangeSort(newSort int) error

	// 删除该单位
	Delete() error
}
