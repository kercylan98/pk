package models

import (
	"pk/common/orm"
	"pk/engine/interfaces"
	"github.com/google/uuid"
	"time"
)

// 排课方案篇章数据结构模型
type PlanSheet struct {
	Id 								string 						`orm:"pk;size(50)" description:"ID"`
	Name 							string						`orm:"size(50)" description:"名称"`
	IsDisable 						bool						`orm:"default(false)" description:"是否失效"`

	Plan	 						*Plan     `orm:"size(50);rel(fk)" description:"所属排课方案"`
	Player 							*Player `orm:"size(50);rel(fk)" description:"创建者玩家"`

	Sort 							int 						`orm:"default(999)" description:"排序序号"`
	TimeCreated						time.Time					`orm:"auto_now_add;null" description:"创建时间"`
	TimeUpdated						time.Time					`orm:"auto_now;null" description:"最后操作时间"`
	orm.Model
}

// 新建一个排课方案篇章
func NewPlanSheet(plan interfaces.PlanInterface, player interfaces.PlayerInterface, sheetName string, sort ...int) *PlanSheet {
	planSheet := &PlanSheet{
		Id:     uuid.New().String(),
		Name:   sheetName,
		Plan:   &Plan{Id: plan.GetId()},
		Player: &Player{Id: player.GetId()},
		Sort:   orm.HasGetOneOrElseDefault(sort, 999),
	}
	return planSheet
}