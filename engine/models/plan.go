package models

import (
	"pk/common/conn"
	"pk/common/orm"
	"pk/engine/interfaces"
	"github.com/google/uuid"
	"time"
)

// 排课计划数据结构模型
// 排课计划应该为总的计划，具体方案内容由各sheet组成，之间规则应该共享
type Plan struct {
	Id 								string 						`orm:"pk;size(50)" description:"ID"`
	Name 							string						`orm:"size(50)" description:"名称"`

	Player	 						*Player `orm:"size(50);rel(fk)" description:"创建者玩家"`

	Sort 							int 						`orm:"default(999)" description:"排序序号"`
	TimeCreated						time.Time					`orm:"auto_now_add;null" description:"创建时间"`
	TimeUpdated						time.Time					`orm:"auto_now;null" description:"最后操作时间"`
	orm.Model
}

func NewPlan(player interfaces.PlayerInterface, planName string, sort ...int) *Plan {
	return &Plan{
		Id:     uuid.New().String(),
		Name:   planName,
		Player: &Player{Id: player.GetId()},
		Sort:   orm.HasGetOneOrElseDefault(sort, 999),
	}
}

func (slf *Plan) IsExist() (bool, error) {
	plan := new(Plan)
	if err := orm.Get().QueryTable("plan").Filter("is_delete", false).
		Filter("player_id", slf.Player.Id).
		Filter("name", slf.Name).Limit(1).One(plan); err != nil {
		if err == orm.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (slf *Plan) GetId() string {
	return slf.Id
}

func (slf *Plan) GetName() string {
	return slf.Name
}

func (slf *Plan) GetPlayer() (interfaces.PlayerInterface, error) {
	err := orm.Get().QueryTable("player").Filter("id", slf.Player.GetId()).Limit(1).One(slf.Player)
	if err != nil {
		return nil, err
	}
	return slf.Player, nil
}

func (slf *Plan) GetSheets(page int, limit int, wheres ...*orm.Where) (*conn.Page, error) {
	qs := orm.Get().QueryTable("plan_sheet").Filter("is_delete", false)
	if count, err := qs.Filter("plan_id", slf.Id).Count(); err != nil {
		return nil, err
	}else {
		for _, where := range wheres {
			qs = where.Handle(qs)
		}
		var sheets []*PlanSheet
		_, err := qs.Limit(orm.CalcPageLimit(page, limit)).All(&sheets)
		if err != nil {
			return nil, err
		}
		return conn.PageUtil(count, page, limit, sheets), nil
	}
}

func (slf *Plan) ChangeName(newName string) error {
	slf.Name = newName
	_, err := orm.Get().Update(slf, "Name")
	return err
}

func (slf *Plan) ChangeSort(newSort int) error {
	slf.Sort = newSort
	_, err := orm.Get().Update(slf, "Sort")
	return err
}

