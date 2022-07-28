package engines

import (
	"errors"
	orm2 "github.com/astaxie/beego/orm"
	"github.com/kercylan98/dev-kits/utils/kstr"
	"pk/common/conn"
	"pk/common/orm"
	"pk/engine/interfaces"
	"pk/engine/models"
	"sync"
)

// 单位引擎
type UnitEngine struct {
	sync.RWMutex
}

func NewUnitEngine() *UnitEngine {
	return &UnitEngine{}
}

// 特定单位是否存在
func (slf *UnitEngine) IsExist(unitId string) (bool, error) {
	if count, err := orm.Get().QueryTable("unit").Filter("is_delete", false).
		Filter("id", unitId).Limit(1).Count(); err != nil {
		return false, err
	} else {
		if count != 0 {
			return true, nil
		} else {
			return false, nil
		}
	}
}

// 创建单位
func (slf *UnitEngine) Create(unitName string, sort ...int) error {
	// 校验单位是否存在
	if count, err := orm.Get().QueryTable("unit").Filter("is_delete", false).
		Filter("name", unitName).Limit(1).Count(); err != nil {
		return err
	} else {
		if count != 0 {
			return errors.New("the unit already exists")
		}
	}
	// 添加单位
	_, err := orm.Get().Insert(models.NewUnit(unitName, sort...))
	return err
}

// 创建多个单位
func (slf *UnitEngine) CreateMulti(sort int, unitNames ...string) error {
	var units []*models.Unit
	for _, unitName := range unitNames {
		units = append(units, models.NewUnit(unitName, sort))
	}
	var findUnits []*models.Unit
	if count, err := orm.Get().QueryTable("unit").Filter("is_delete", false).Filter("name__in", unitNames).All(&findUnits, "name"); err != nil {
		return err
	} else if count > 0 {
		source := "these units already exist: "
		for _, unit := range findUnits {
			source += unit.Name + ", "
		}
		return errors.New(kstr.RemoveLast(kstr.RemoveLast(source)))
	} else {
		rollbackOrm := orm2.NewOrm()
		err := rollbackOrm.Begin()
		if err != nil {
			return err
		}
		_, err = rollbackOrm.InsertMulti(100, units)
		if err != nil {
			rollbackErr := rollbackOrm.Rollback()
			if rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
		return rollbackOrm.Commit()
	}
}

// 获取特定单位信息
func (slf *UnitEngine) Get(unitId string) (interfaces.UnitInterface, error) {
	unit := new(models.Unit)
	err := orm.Get().QueryTable("unit").Filter("is_delete", false).Filter("id", unitId).Limit(1).One(unit)
	return unit, err
}

// 获取多个单位
func (slf *UnitEngine) GetMulti(page, limit int, wheres ...*orm.Where) (*conn.Page, error) {
	qt := orm.Get().QueryTable("unit").Filter("is_delete", false)
	if count, err := qt.Count(); err != nil {
		return nil, err
	} else {
		var units []*models.Unit
		for _, where := range wheres {
			qt = where.Handle(qt)
		}
		_, err := qt.Limit(orm.CalcPageLimit(page, limit)).OrderBy("sort", "-time_created").All(&units)
		return conn.PageUtil(count, page, limit, units), err
	}
}
