package models

import (
	"pk/common/orm"
	orm2 "github.com/astaxie/beego/orm"
	"github.com/google/uuid"
	"time"
)

// 排课方案工作日数据结构模型
type PlanWorkday struct {
	Id 								string 						`orm:"pk;size(50)" description:"ID"`
	Name 							string						`orm:"size(50)" description:"名称"`
	Week 							int 						`orm:"size(1)" description:"周次"`
	IsDisable 						bool						`orm:"default(false)" description:"是否失效"`

	SheetId	 						string 						`orm:"size(50)" description:"所属排课方案篇章ID"`

	TimeCreated						time.Time					`orm:"auto_now_add;null" description:"创建时间"`
	TimeUpdated						time.Time					`orm:"auto_now;null" description:"最后操作时间"`
	orm.Model
}

// 新建一个工作日
// 在创建工作日的时候应该将其他相同工作日进行停用
func NewPlanWorkday(sheetId, workdayName string, week int) (*PlanWorkday, error) {
	workday := &PlanWorkday{
		Id:   uuid.New().String(),
		Name: workdayName,
		SheetId: sheetId,
		Week: week,
	}
	tempOrm := orm2.NewOrm()
	if err := tempOrm.Begin(); err != nil {
		return nil, err
	}
	// 设置其他相同篇章且相同工作日的数据停用，若是出现异常则回滚返回错误
	if _, err := tempOrm.QueryTable(workday).
		Filter("SheetId", sheetId).
		Filter("Week", week).
		Update(orm2.Params{"is_disable": true}); err != nil {
		if rollbackErr := tempOrm.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}
	// 尝试插入新的工作日数据，若是出现异常则回滚之前的设置操作并返回错误
	if _, err := tempOrm.Insert(workday); err != nil {
		if rollbackErr := tempOrm.Rollback(); rollbackErr != nil {
			return nil, rollbackErr
		}
		return nil, err
	}
	// 尝试提交
	if err := tempOrm.Commit(); err != nil {
		return nil, err
	}

	return workday, nil
}