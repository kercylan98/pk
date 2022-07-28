package models

import "pk/common/orm"

// 单位班级教学学科数据结构
type UnitClassCourseTeacherR struct {
	Id 					string      `orm:"pk;size(50)" description:"ID"`
	UnitClass			*UnitClass   `orm:"index;rel(one);on_delete(do_nothing)" description:"班级"`
	UnitCourse 			*UnitCourse `orm:"index;rel(one);on_delete(do_nothing);null" description:"课程"`
	UnitTeacher 		*UnitTeacher   `orm:"index;rel(one);on_delete(do_nothing);null" description:"课程任教老师"`
	WeekNumber 			int         `orm:"size(2);default(0)" description:"周课时数"`
	orm.Model

}
