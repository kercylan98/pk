package models

import (
	"pk/common/orm"
	"pk/engine/interfaces"
	"github.com/google/uuid"
	"time"
)

// 教师结构定义
type UnitTeacher struct {
	Id 					string					`orm:"pk;size(50)" description:"ID"`
	UnitId 				string					`orm:"size(50)" description:"所属单位ID"`
	Name 				string					`orm:"size(50)" description:"名称"`

	Sort 				int 					`orm:"default(9999)" description:"排序序号"`
	TimeCreated			time.Time				`orm:"auto_now_add;null" description:"创建时间"`
	TimeUpdated			time.Time				`orm:"auto_now;null" description:"最后操作时间"`
	orm.Model

}

func NewUnitTeacher(unit interfaces.UnitInterface, teacherName string, sort ...int) *UnitTeacher {
	return &UnitTeacher{
		Id:     uuid.New().String(),
		UnitId: unit.GetId(),
		Name:   teacherName,
		Sort:   orm.HasGetOneOrElseDefault(sort, 9999),
	}
}

func (slf *UnitTeacher) GetId() string {
	return slf.Id
}

func (slf *UnitTeacher) GetName() string {
	panic("implement me")
}

func (slf *UnitTeacher) ChangeName(newName string) error {
	panic("implement me")
}

func (slf *UnitTeacher) ChangeSort(newSort int) error {
	panic("implement me")
}

func (slf *UnitTeacher) Delete() error {
	panic("implement me")
}

