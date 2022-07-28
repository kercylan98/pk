package models

import (
	"pk/common/conn"
	"pk/common/orm"
	"pk/engine/interfaces"
	"github.com/google/uuid"
	"time"
)

// 单位班级数据结构
type UnitClass struct {
	Id 					string    `orm:"pk;size(50)" description:"ID"`
	Grade 				*UnitGrade `orm:"size(50);rel(fk)" description:"所属年级"`
	Teacher 			*UnitTeacher `orm:"size(50);null;rel(fk)" description:"班主任教师"`
	Name 				string      `orm:"size(50)" description:"名称"`

	Sort 				int 					`orm:"default(999)" description:"排序序号"`
	TimeCreated			time.Time				`orm:"auto_now_add;null" description:"创建时间"`
	TimeUpdated			time.Time				`orm:"auto_now;null" description:"最后操作时间"`
	orm.Model

}

func NewUnitClass(unitGrade interfaces.UnitGradeInterface, className string, headTeacher interfaces.UnitTeacherInterface, sort ...int) *UnitClass {
	return &UnitClass{
		Id:      uuid.New().String(),
		Grade:   &UnitGrade{Id: unitGrade.GetId()},
		Teacher: &UnitTeacher{Id: headTeacher.GetId()},
		Name:    className,
		Sort:    orm.HasGetOneOrElseDefault(sort, 999),
	}
}

func NewUnitClassNoTeacher(unitGrade interfaces.UnitGradeInterface, className string, sort ...int) *UnitClass {
	return &UnitClass{
		Id:    uuid.New().String(),
		Grade: &UnitGrade{Id: unitGrade.GetId()},
		Name:  className,
		Sort:  orm.HasGetOneOrElseDefault(sort, 999),
	}
}

func (slf *UnitClass) GetId() string {
	return slf.Id
}

func (slf *UnitClass) GetName() string {
	return slf.Name
}

func (slf *UnitClass) GetGrade() (interfaces.UnitGradeInterface, error) {
	unitGrade := new(UnitGrade)
	err := orm.Get().QueryTable(unitGrade).Filter("id", slf.Grade.Id).Limit(1).One(unitGrade)
	return unitGrade, err
}

func (slf *UnitClass) GetHeadTeacher() (interfaces.UnitTeacherInterface, error) {
	headTeacher := new(UnitTeacher)
	err := orm.Get().QueryTable(headTeacher).Filter("id", slf.Teacher.Id).Limit(1).One(headTeacher)
	return headTeacher, err
}

func (slf *UnitClass) GetTeachers(page, limit int, wheres ...*orm.Where) (*conn.Page, error) {
	qs := orm.Get().QueryTable("unit_class_course_teacher_r").Filter("is_delete", false)
	qs.RelatedSel("unit_teacher").Filter("UnitTeacher__is_delete", false).Distinct()
	if count, err := qs.Count(); err != nil {
		return nil, err
	}else {
		for _, where := range wheres {
			qs = where.Handle(qs)
		}
		var teachers []*UnitTeacher
		_, err := qs.Limit(orm.CalcPageLimit(page, limit)).All(&teachers, "UnitTeacher")
		return conn.PageUtil(count, page, limit, teachers), err
	}
}

func (slf *UnitClass) GetCourses(page, limit int, wheres ...*orm.Where) (*conn.Page, error) {
	qs := orm.Get().QueryTable("unit_class_course_teacher_r").Filter("is_delete", false)
	qs.RelatedSel("unit_course").Filter("UnitCourse__is_delete", false).Distinct()
	if count, err := qs.Count(); err != nil {
		return nil, err
	}else {
		for _, where := range wheres {
			qs = where.Handle(qs)
		}
		var courses []*UnitCourse
		_, err := qs.Limit(orm.CalcPageLimit(page, limit)).All(&courses, "UnitCourse")
		return conn.PageUtil(count, page, limit, courses), err
	}
}

func (slf *UnitClass) AddCourse(course interfaces.UnitCourseInterface, weekNumber int, teachers ...interfaces.UnitTeacherInterface) {
	panic("implement me")
}

func (slf *UnitClass) SetWeekNumber(course interfaces.UnitCourseInterface, weekNumber int) error {
	panic("implement me")
}

func (slf *UnitClass) SetTeacher(course interfaces.UnitCourseInterface, teachers ...interfaces.UnitTeacherInterface) {
	panic("implement me")
}

func (slf *UnitClass) ChangeHeadTeacher(teacher interfaces.UnitTeacherInterface) error {
	panic("implement me")
}

func (slf *UnitClass) ChangeName(newName string) error {
	panic("implement me")
}

func (slf *UnitClass) ChangeSort(newSort int) error {
	panic("implement me")
}

func (slf *UnitClass) Delete() error {
	panic("implement me")
}

func (slf *UnitClass) DeleteCourse(course interfaces.UnitCourseInterface) error {
	panic("implement me")
}

func (slf *UnitClass) DeleteTeacher(course interfaces.UnitCourseInterface, teacher interfaces.UnitTeacherInterface) error {
	panic("implement me")
}

