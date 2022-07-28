package engines

import (
	"pk/common/orm"
	"pk/engine/interfaces"
	"pk/engine/models"
	"errors"
)

type PlanEngine struct {

}

func NewPlanEngine() *PlanEngine {
	return &PlanEngine{}
}

// 创建排课计划
func (slf *PlanEngine) Create(player interfaces.PlayerInterface, name string, sort ...int) error {
	plan := models.NewPlan(player, name, sort...)
	if exist, err := plan.IsExist(); err != nil {
		return err
	}else if exist {
		return errors.New("the plan already exists")
	}
	_, err := orm.Get().Insert(plan)
	if err != nil {
		return err
	}
	return nil
}