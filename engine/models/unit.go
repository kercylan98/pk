package models

import (
	"pk/common/conn"
	"pk/common/orm"
	"pk/engine/interfaces"
	"errors"
	"github.com/google/uuid"
	"time"
)

// 单位数据结构模型
type Unit struct {
	Id 								string 						`orm:"pk;size(50)" description:"ID"`
	Name 							string						`orm:"size(50);unique" description:"名称"`
	SynchronousSource				string						`orm:"size(200);null" description:"同步数据来源URL"`

	Sort 							int 						`orm:"default(999)" description:"排序序号"`
	TimeCreated						time.Time					`orm:"auto_now_add;null" description:"创建时间"`
	TimeUpdated						time.Time					`orm:"auto_now;null" description:"最后操作时间"`
	orm.Model

}

func NewUnit(unitName string, sort ...int) *Unit {
	return &Unit{
		Id:   uuid.New().String(),
		Name: unitName,
		Sort: orm.HasGetOneOrElseDefault(sort, 0),
	}
}

func (slf *Unit) GetId() string {
	return slf.Id
}

func (slf *Unit) Delete() error {
	return slf.Del()
}

func (slf *Unit) GetName() string {
	return slf.Name
}

func (slf *Unit) GetPlayers(page, limit int, wheres ...*orm.Where) (*conn.Page, error) {
	qt := orm.Get().QueryTable("player").Filter("unit_id", slf.Id).Filter("is_delete", false)
	if count, err := qt.Count(); err != nil {
		return nil, err
	}else {
		for _, where := range wheres {
			qt = where.Handle(qt)
		}
		var players []*Player
		_, err := qt.Limit(orm.CalcPageLimit(page, limit)).All(&players)
		return conn.PageUtil(count, page, limit, players), err
	}
}

func (slf *Unit) GetGrades(page, limit int, wheres ...*orm.Where) (*conn.Page, error) {
	qt := orm.Get().QueryTable("unit_grade").Filter("unit_id", slf.Id).Filter("is_delete", false)
	if count, err := qt.Count(); err != nil {
		return nil, err
	}else {
		for _, where := range wheres {
			qt = where.Handle(qt)
		}
		var unitGrades []*UnitGrade
		_, err := qt.Limit(orm.CalcPageLimit(page, limit)).All(&unitGrades)
		return conn.PageUtil(count, page, limit, unitGrades), err
	}
}

func (slf *Unit) GetSites(page, limit int, wheres ...*orm.Where) (*conn.Page, error) {
	qt := orm.Get().QueryTable("unit_site").Filter("unit_id", slf.Id).Filter("is_delete", false)
	if count, err := qt.Count(); err != nil {
		return nil, err
	}else {
		for _, where := range wheres {
			qt = where.Handle(qt)
		}
		var unitSites []*UnitSite
		_, err := qt.Limit(orm.CalcPageLimit(page, limit)).All(&unitSites)
		return conn.PageUtil(count, page, limit, unitSites), err
	}
}

func (slf *Unit) GetTeachers(page, limit int, wheres ...*orm.Where) (*conn.Page, error) {
	qt := orm.Get().QueryTable("unit_teacher").Filter("unit_id", slf.Id).Filter("is_delete", false)
	if count, err := qt.Count(); err != nil {
		return nil, err
	}else {
		for _, where := range wheres {
			qt = where.Handle(qt)
		}
		var unitTeachers []*UnitTeacher
		_, err := qt.Limit(orm.CalcPageLimit(page, limit)).All(&unitTeachers)
		return conn.PageUtil(count, page, limit, unitTeachers), err
	}
}

func (slf *Unit) AddGrade(gradeName string, sort ...int) (interfaces.UnitGradeInterface, error) {
	unitGrade := NewUnitGrade(slf, gradeName, sort...)
	if count, err := orm.Get().QueryTable(unitGrade).Filter("name", gradeName).Count(); err != nil {
		return nil, err
	}else if count > 0 {
		return nil, errors.New("the grade already exists and is not allowed to be created repeatedly")
	}
	_, err := orm.Get().Insert(unitGrade)
	return unitGrade, err
}

func (slf *Unit) AddSite(siteName string, sort ...int) (interfaces.UnitSiteInterface, error) {
	unitSite := NewUnitSite(slf, siteName, sort...)
	if count, err := orm.Get().QueryTable(unitSite).Filter("name", siteName).Count(); err != nil {
		return nil, err
	}else if count > 0 {
		return nil, errors.New("the site already exists and is not allowed to be created repeatedly")
	}
	_, err := orm.Get().Insert(unitSite)
	return unitSite, err
}

func (slf *Unit) AddTeacher(teacherName string, sort ...int) (interfaces.UnitTeacherInterface, error) {
	unitTeacher := NewUnitSite(slf, teacherName, sort...)
	if count, err := orm.Get().QueryTable(unitTeacher).Filter("name", teacherName).Count(); err != nil {
		return nil, err
	}else if count > 0 {
		return nil, errors.New("the teacher already exists and is not allowed to be created repeatedly")
	}
	_, err := orm.Get().Insert(unitTeacher)
	return unitTeacher, err
}

func (slf *Unit) ChangeName(newName string) error {
	slf.Name = newName
	_, err := orm.Get().Update(slf, "Name")
	return err
}

func (slf *Unit) ChangeSort(newSort int) error {
	slf.Sort = newSort
	_, err := orm.Get().Update(slf, "Sort")
	return err
}
