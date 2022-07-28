package controllers

import (
	"pk/common/basic"
	"pk/common/conn"
	"pk/engine/checks"
	"pk/engine/models"
	"fmt"
	"strings"
)

type PlayerController struct {
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
func Get()  {

}


// @Title 创建玩家
// @Description 创建一个或多个玩家，当存在任意一个玩家创建失败时，所有玩家均视为创建失败
// @Param	uid	form	string	true	"需要创建玩家的单位id"
// @Param	account	form	string	false	"需要创建的玩家帐号，当需要同时创建多个的时候使用英文逗号分隔"
// @Param	name	from	string	false	"为创建的玩家设置名字，当创建帐号为多个的时候该参数无效，名称默认同帐号相同"
// @Param	pwd	from	string	false	"为创建的玩家设置密码，当创建帐号为多个的时候该参数无效，密码默认同帐号相同"
// @Param	sort	form	int	false	"为创建的玩家同时设置序号,若未设置则为默认999"
// @Param	load	form	string	false	"通过文件模板进行玩家的批量导入时需要设置该参数为文件上传的name，存在该参数时，account、pwd、name和sort失效"
// @Success successful:0 {object} conn.ResponseBodyModel
// @Failure unexpected:-1 系统错误
// @Failure unexpected:40002 未找到特定单位
// @Failure unexpected:40003 不合法的请求参数
// @Failure unexpected:50000 待创建的账户已存在
// @router / [post]
func (slf *PlayerController) Post() {
	slf.Data["json"] = conn.Dispose(func(unitId, accounts, password, name, load string) (interface{}, interface{}) {
		sort, _ := slf.GetInt("sort", 999)
		if exist, err := slf.App().UnitEngine.IsExist(unitId); err != nil {
			return conn.CODE_SYSTEM_ERR, "err:" + err.Error()
		}else if !exist {
			return conn.CODE_NOTFOUND, "err:尝试向一个不存在的单位中创建玩家，创建失败"
		}

		unit := &models.Unit{Id: unitId}
		// 发生错误则尝试通过表单数据创建玩家

		if _, _, err := slf.GetFile("load"); err != nil {
			// 解析帐号
			var playerAccounts []string
			for _, account := range strings.Split(accounts, ",") {
				// 发生错误时尝试通过参数进行创建，否则通过上传的文件进行
				if v1, v2, ok := checks.CheckPlayerAccount(account); !ok {
					return v1, v2
				}
				playerAccounts = append(playerAccounts, account)
			}
			// 判断单玩家创建还是多玩家创建
			if len(playerAccounts) > 1 {
				if err := slf.App().PlayerEngine.RegistryMulti(unit, sort, playerAccounts...); err != nil {
					key := "these players already exist:"
					if strings.Contains(err.Error(), key) {
						return conn.CODE_EXIST, fmt.Sprintf("err:批量创建账户失败，这些账户已经存在：%s", strings.ReplaceAll(err.Error(), key, ""))
					}
				}
				return nil, nil
			} else {
				if v1, v2, ok := checks.CheckPlayerPassword(password); !ok {
					return v1, v2
				}
				if v1, v2, ok := checks.CheckPlayerName(name); !ok {
					return v1, v2
				}
				err := slf.App().PlayerEngine.Registry(
					unit, name, strings.TrimSpace(accounts), password, sort)
				if err != nil {
					if err.Error() == "the units don't exist" {
						return conn.CODE_NOTFOUND, "err:尝试向一个不存在的单位中创建玩家，创建失败"
					} else if err.Error() == "account already in existence" {
						return conn.CODE_EXIST, "err:该账户已存在"
					}
					return conn.CODE_SYSTEM_ERR, err
				}
			}
			return nil, nil
		}else {
			if files, _, err := slf.App().FileEngine.Dispose(slf.Ctx.Request); err != nil {
				return conn.CODE_SYSTEM_ERR, "err:" + err.Error()
			}else {
				for _, file := range files {
					if data, err := slf.App().XlsxEngine.GetAll(file.GetSavePath(), 0, 1); err == nil {
						var players []*models.Player
						for i, line := range data {
							// 处理数据
							account := line[2].String()
							if account == "" {
								return conn.CODE_ILLEGAL_ARGS, fmt.Sprint("err:导入模板第", i, "行数据帐号不能为空")
							}
							name := line[1].String()
							password := line[3].String()
							sort, err := line[0].Int()
							if err != nil {
								sort = 999
							}
							if name == "" {
								name = account
							}
							if password == "" {
								password = account
							}
							players = append(players, models.NewPlayer(unit, name, account, password, sort))
						}
						// 入库
						err := slf.App().PlayerEngine.RegistryMultiWithModel(unit, players...)
						if err != nil {
							if strings.Contains(err.Error(), "these players already exist") {
								return conn.CODE_EXIST, "err:" + err.Error()
							}
							return conn.CODE_SYSTEM_ERR, err
						}
					}else {
						if strings.Contains(err.Error(), "zip: not a valid zip file") {
							return conn.CODE_SYSTEM_ERR, "err:即将导入的模板文件不是一个有效的xlsx文件"
						}
						return conn.CODE_SYSTEM_ERR, "err:" + err.Error()
					}
				}
			}

			return nil, nil
		}
	},
		slf.GetString, "uid",
		slf.GetString, "account",
		slf.GetString, "pwd",
		slf.GetString, "name",
		slf.GetString, "load")
	slf.ServeJSON()
}