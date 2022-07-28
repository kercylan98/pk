package interfaces

import (
	"pk/common/conn"
	"pk/common/orm"
)

type UnitClassInterface interface {
	// 获取ID
	GetId() string
	// 获取名称
	GetName() string

	// 获取年级信息
	GetGrade() (UnitGradeInterface, error)
	// 获取班主任信息
	GetHeadTeacher() (UnitTeacherInterface, error)
	// 获取班级下所有任课老师信息
	GetTeachers(page, limit int, wheres ...*orm.Where) (*conn.Page, error)
	// 获取班级下所有学科信息
	GetCourses(page, limit int, wheres ...*orm.Where) (*conn.Page, error)

	// 添加班级课程
	AddCourse(course UnitCourseInterface, weekNumber int, teachers ...UnitTeacherInterface)

	// 设置班级课程周课时数
	SetWeekNumber(course UnitCourseInterface, weekNumber int) error
	// 设置班级课程任课教师
	SetTeacher(course UnitCourseInterface, teachers ...UnitTeacherInterface)

	// 改变班主任
	ChangeHeadTeacher(teacher UnitTeacherInterface) error
	// 改变名称
	ChangeName(newName string) error
	// 改变序号
	ChangeSort(newSort int) error

	// 删除班级
	Delete() error
	// 删除班级课程
	DeleteCourse(course UnitCourseInterface) error
	// 删除班级课程任课教师
	DeleteTeacher(course UnitCourseInterface, teacher UnitTeacherInterface) error
}
