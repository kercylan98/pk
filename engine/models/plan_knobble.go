package models

import (
	"pk/common/orm"
	orm2 "github.com/astaxie/beego/orm"
	"github.com/google/uuid"
	"time"
)

// 排课方案课节次数据结构模型
type PlanKnobble struct {
	Id 								string 						`orm:"pk;size(50)" description:"ID"`
	Name 							string						`orm:"size(50)" description:"名称"`
	Section 						int 						`orm:"size(2)" description:"上课节次"`
	StartTime 						time.Time					`description:"课节次开始时间"`
	EndTime 						time.Time					`description:"课节次结束时间"`
	IsDisable 						bool						`orm:"default(false)" description:"是否失效"`

	WorkdayId	 					string 						`orm:"size(50)" description:"所属工作日ID"`

	TimeCreated						time.Time					`orm:"auto_now_add;null" description:"创建时间"`
	TimeUpdated						time.Time					`orm:"auto_now;null" description:"最后操作时间"`
	orm.Model
}

// 新建课节次
// 在创建课节次的时候应该将其他相同课节次进行停用
func NewPlanKnobble(workdayId, knobbleName string, section int) (*PlanKnobble, error) {
	knobble := &PlanKnobble{
		Id:   uuid.New().String(),
		Name: knobbleName,
		WorkdayId: workdayId,
		Section: section,
	}
	tempOrm := orm2.NewOrm()
	if err := tempOrm.Begin(); err != nil {
		return nil, err
	}
	// 设置其他相同课节次数据停用，若是出现异常则回滚返回错误
	if _, err := tempOrm.QueryTable(knobble).
		Filter("WorkdayId", workdayId).
		Filter("Section", section).
		Update(orm2.Params{"is_disable": true}); err != nil {
		if rollbackErr := tempOrm.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}
	// 尝试插入新的课节次数据，若是出现异常则回滚之前的设置操作并返回错误
	if _, err := tempOrm.Insert(knobble); err != nil {
		if rollbackErr := tempOrm.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}
	// 尝试提交
	if err := tempOrm.Commit(); err != nil {
		return nil, err
	}

	return knobble, nil
}