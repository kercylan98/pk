package controllers

import (
	"pk/common/basic"
	"pk/common/conn"
	"pk/common/orm"
	"pk/common/orm/wheres"
	"pk/engine/checks"
	"fmt"
	"strings"
)

type UnitController struct {
	basic.Controller
}

// @Title 获取单位信息
// @Description 获取单位信息，当未明确获取目标时则为可附加条件的单位列表获取
// @Param	uid	path	string	false	"需要获取信息的单位ID，如果无该参数则表示获取列表"
// @Param	p	query	int	false	"需要获取的页码，默认为1，仅在获取列表时生效"
// @Param	l	query	int	false	"需要显示的条数，默认为10，仅在获取列表时生效"
// @Param	n	query	string	false	"需要筛选的单位名称，仅在获取列表时生效"
// @Success successful:0 {object} conn.ResponseBodyModel
// @Success successful:0-A {object} models.Unit
// @Success successful:0-B {object} []models.Unit
// @Failure unexpected:-1 系统错误
// @Failure unexpected:4ooo2 尝试访问一个不存在的单位，获取单位信息失败
// @router /?:uid [get]
func (slf *UnitController) Get() {
	slf.Data["json"] = conn.Dispose(func(unitId, containUnitName string) (interface{}, interface{}) {
		page, _ := slf.GetInt("p", 1)
		limit, _ := slf.GetInt("l", 10)
		if len(strings.TrimSpace(unitId)) == 0 {
			if page, err := slf.App().UnitEngine.GetMulti(page, limit, new(orm.Where).
				Add(wheres.Contains("name", strings.TrimSpace(containUnitName))),
			); err != nil {
				return conn.CODE_SYSTEM_ERR, err
			}else {
				return page, nil
			}
		}else {
			unit, err := slf.App().UnitEngine.Get(unitId)
			if err != nil {
				if orm.EqErr(err, orm.ErrNoRows) {
					return conn.CODE_NOTFOUND, fmt.Sprintf("err:尝试访问一个不存在的单位(%s)，获取单位信息失败", unitId)
				}
				return conn.CODE_SYSTEM_ERR, fmt.Sprint("err:", err)
			}
			return unit, nil
		}

	},
		slf.GetString, ":uid",
		slf.GetString, "n",)
	slf.ServeJSON()
}

// @Title 创建单位
// @Description 创建一个或多个单位，当存在任意一个单位创建失败时，所有单位均视为创建失败
// @Param	name	form	string	true	"需要创建的单位名称，当需要同时创建多个单位的时候使用英文逗号分隔"
// @Param	sort	form	int	false	"为创建的单位同时设置序号,若未设置则为默认999"
// @Success successful:0 {object} conn.ResponseBodyModel
// @Failure unexpected:-1 系统错误
// @Failure unexpected:50000 数据已存在
// @Failure unexpected:40003 不合法的请求参数
// @router / [post]
func (slf *UnitController) Post() {
	slf.Data["json"] = conn.Dispose(func(unitNames string) (interface{}, interface{}) {
		sort, _ := slf.GetInt("sort", 999)
		names := strings.Split(unitNames, ",")
		// 格式化单位名称
		var formatNames []string
		for _, name := range names {
			formatName := strings.TrimSpace(name)
			if v1, v2, ok := checks.CheckUnitName(formatName); !ok {
				return v1, v2
			}
			formatNames = append(formatNames, formatName)
		}
		if len(names) == 1 {
			if err := slf.App().UnitEngine.Create(names[0], sort); err != nil {
				if strings.Contains(err.Error(), "the unit already exists") {
					return conn.CODE_EXIST, "err:该单位已存在"
				}
				return conn.CODE_SYSTEM_ERR, "err:" + err.Error()
			}
			return nil, nil
		}else {
			if err := slf.App().UnitEngine.CreateMulti(sort, formatNames...); err != nil {
				key := "these units already exist:"
				if strings.Contains(err.Error(), key) {
					return conn.CODE_EXIST, fmt.Sprintf("err:批量创建单位失败，这些单位已经存在：%s", strings.ReplaceAll(err.Error(), key, ""))
				}
			}
			return nil, nil
		}
	},
		slf.GetString, "name",)
	slf.ServeJSON()
}