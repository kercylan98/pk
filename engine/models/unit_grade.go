package models

import (
	"pk/common/conn"
	"pk/common/orm"
	"pk/engine/interfaces"
	"github.com/google/uuid"
	"time"
)

// 单位年级数据结构
type UnitGrade struct {
	Id 					string					`orm:"pk;size(50)" description:"ID"`
	UnitId 				string					`orm:"size(50)" description:"所属单位ID"`
	Name 				string					`orm:"size(50)" description:"名称"`

	Sort 				int 					`orm:"default(999)" description:"排序序号"`
	TimeCreated			time.Time				`orm:"auto_now_add;null" description:"创建时间"`
	TimeUpdated			time.Time				`orm:"auto_now;null" description:"最后操作时间"`
	orm.Model

}

func NewUnitGrade(unit interfaces.UnitInterface, gradeName string, sort ...int) *UnitGrade {
	s := 999
	if len(sort) > 0 {
		s = sort[0]
	}
	return &UnitGrade{
		Id:          uuid.New().String(),
		UnitId:      unit.GetId(),
		Name:        gradeName,
		Sort:        s,
	}
}

func (slf *UnitGrade) GetId() string {
	return slf.Id
}

func (slf *UnitGrade) GetName() string {
	panic("implement me")
}

func (slf *UnitGrade) GetUnit() (interfaces.UnitInterface, error) {
	panic("implement me")
}

func (slf *UnitGrade) GetClasses(page, limit int, wheres ...*orm.Where) (*conn.Page, error) {
	panic("implement me")
}

func (slf *UnitGrade) ChangeName(newName string) error {
	panic("implement me")
}

func (slf *UnitGrade) ChangeSort(newSort int) error {
	panic("implement me")
}

