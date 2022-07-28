package models

import (
	"pk/common/orm"
	"pk/engine/interfaces"
	"github.com/google/uuid"
	"time"
)

// 单位教学学科数据结构
type UnitCourse struct {
	Id 					string					`orm:"pk;size(50)" description:"ID"`
	UnitId 				string					`orm:"size(50)" description:"所属单位ID"`
	Name 				string					`orm:"size(50)" description:"名称"`

	Sort 				int 					`orm:"default(9999)" description:"排序序号"`
	TimeCreated			time.Time				`orm:"auto_now_add;null" description:"创建时间"`
	TimeUpdated			time.Time				`orm:"auto_now;null" description:"最后操作时间"`
	orm.Model

}

func NewUnitCourse(unit interfaces.UnitInterface, courseName string, sort ...int) *UnitCourse {
	return &UnitCourse{
		Id:     uuid.New().String(),
		UnitId: unit.GetId(),
		Name:   courseName,
		Sort:   orm.HasGetOneOrElseDefault(sort, 9999),
	}
}

func (slf *UnitCourse) GetId() string {
	return slf.Id
}

func (slf *UnitCourse) GetName() string {
	panic("implement me")
}

func (slf *UnitCourse) ChangeName(newName string) error {
	panic("implement me")
}

func (slf *UnitCourse) ChangeSort(newSort int) error {
	panic("implement me")
}

func (slf *UnitCourse) Delete() error {
	panic("implement me")
}

