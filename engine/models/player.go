package models

import (
	"pk/common/conn"
	"pk/common/orm"
	"pk/engine/interfaces"
	"errors"
	"github.com/google/uuid"
	"time"
)

// 实际排课方案操作者（玩家）数据模型
type Player struct {
	Id 					string  `orm:"pk;size(50)" description:"ID"`
	Unit 				*Unit     `orm:"size(50);rel(fk)" description:"所属单位"`
	Name 				string    `orm:"size(50)" description:"名称"`
	Account				string `orm:"size(50)" description:"登录帐号"`
	Password 			string    `orm:"size(100)" description:"登录密码" json:"-"`
	IsSynchronous		bool     `orm:"default(false)" description:"是否是同步账户，如果是应该单位同步源登录"`
	LastOperationId		string `orm:"size(50);null" description:"最后操作的方案ID"`

	Sort 				int 					`orm:"default(999)" description:"排序序号"`
	TimeLoggedIn		time.Time				`orm:"null" description:"最后登录时间"`
	TimeCreated			time.Time				`orm:"auto_now_add;null" description:"注册时间"`
	TimeUpdated			time.Time				`orm:"auto_now;null" description:"最后操作时间"`
	orm.Model

}

func NewPlayer(unit interfaces.UnitInterface, playerName, playerAccount, playerPassword string, sort ...int) *Player {
	return &Player{
		Id:       uuid.New().String(),
		Unit:     &Unit{Id: unit.GetId()},
		Name:     playerName,
		Account:  playerAccount,
		Password: playerPassword,
		Sort:     orm.HasGetOneOrElseDefault(sort, 999),
	}
}

func (slf *Player) GetId() string {
	return slf.Id
}

func (slf *Player) ChangeSort(newSort int) error {
	slf.Sort = newSort
	_, err := orm.Get().Update(slf, "Sort")
	return err
}

func (slf *Player) GetPlans(page, limit int, wheres ...*orm.Where) (*conn.Page, error) {
	qs := orm.Get().QueryTable("plan").Filter("is_delete", false)
	if count, err := qs.Count(); err != nil {
		return nil, err
	}else {
		for _, where := range wheres {
			qs = where.Handle(qs)
		}
		var plans []*Plan
		_, err := qs.Limit(orm.CalcPageLimit(page, limit)).All(&plans)
		if err != nil {
			return nil, err
		}
		return conn.PageUtil(count, page, limit, plans), nil
	}
}

func (slf *Player) CreatePlan(planName string, sort ...int) (interfaces.PlanInterface, error) {
	plan := NewPlan(slf, planName, sort...)
	err := orm.Get().QueryTable("plan").Filter("name", planName).Filter("player_id", slf.Id).Limit(1).One(new(Plan))
	if err != nil {
		if err == orm.ErrNoRows {
			_, err = orm.Get().Insert(plan)
			return plan, err
		}else {
			return nil, err
		}
	}else {
		return nil, errors.New("plan existed, each person can only create one scheme with the same name")
	}
}

func (slf *Player) ChangePassword(newPassword string) error {
	slf.Password = newPassword
	_, err := orm.Get().Update(slf, "Password")
	return err
}

func (slf *Player) GetName() string {
	return slf.Name
}

func (slf *Player) ChangeName(newName string) error {
	slf.Name = newName
	_, err := orm.Get().Update(slf, "Name")
	return err
}

func (slf *Player) GetUnit() (interfaces.UnitInterface, error) {
	unit := new(Unit)
	err := orm.Get().QueryTable("unit").Filter("id", slf.Unit.Id).Limit(1).One(unit)
	return unit, err
}