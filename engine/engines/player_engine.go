package engines

import (
	"errors"
	"github.com/astaxie/beego/logs"
	orm2 "github.com/astaxie/beego/orm"
	"github.com/dgrijalva/jwt-go"
	"github.com/kercylan98/dev-kits/utils/kstr"
	"pk/common/orm"
	"pk/engine/interfaces"
	"pk/engine/models"
	"strings"
	"sync"
	"time"
)

const SECRET_KEY = "SCHOOL_SCHEDULE_PLAN_DESIGNER_KERCYLAN."

// 玩家引擎
type PlayerEngine struct {
	sync.RWMutex
}

func NewPlayerEngine() *PlayerEngine {
	return &PlayerEngine{}
}

// 注册多个玩家到排课系统
func (slf *PlayerEngine) RegistryMultiWithModel(unit interfaces.UnitInterface, players ...*models.Player) error {
	// 校验单位是否存在
	if count, err := orm.Get().QueryTable("unit").Filter("is_delete", false).
		Filter("id", unit.GetId()).Count(); err != nil {
		return err
	} else if count <= 0 {
		return errors.New("the units don't exist")
	}

	// 检查用户是否存在
	var formatAccounts []string
	for _, player := range players {
		formatAccounts = append(formatAccounts, player.Account)
	}
	var findPlayers []*models.Player
	if count, err := orm.Get().QueryTable("player").
		Filter("is_delete", false).
		Filter("unit_id", unit.GetId()).
		Filter("account__in", formatAccounts).
		All(&findPlayers, "account"); err != nil {
		return err
	} else if count > 0 {
		source := "these players already exist: "
		for _, player := range findPlayers {
			source += player.Account + ", "
		}
		return errors.New(kstr.RemoveLast(kstr.RemoveLast(source)))
	} else {
		// 入库
		rollbackOrm := orm2.NewOrm()
		err := rollbackOrm.Begin()
		if err != nil {
			return err
		}
		_, err = rollbackOrm.InsertMulti(100, players)
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

// 注册一个玩家到整个排课系统，返回错误信息
func (slf *PlayerEngine) Registry(unit interfaces.UnitInterface, playerName, playerAccount, playerPassword string, sort ...int) error {
	// 校验单位是否存在
	if count, err := orm.Get().QueryTable("unit").Filter("is_delete", false).
		Filter("id", unit.GetId()).Count(); err != nil {
		return err
	} else if count <= 0 {
		return errors.New("the units don't exist")
	}

	// 校验用户是否存在
	if err := orm.Get().QueryTable("player").
		Filter("unit_id", unit.GetId()).
		Filter("account", playerAccount).Limit(1).One(new(models.Player)); err == nil {
		return errors.New("account already in existence")
	} else {
		if orm.EqErr(err, orm.ErrNoRows) {
			goto reg
		} else {
			return err
		}
	}

reg:
	{
		// 插入用户数据
		player := models.NewPlayer(unit, playerName, playerAccount, playerPassword, sort...)
		if _, err := orm.Get().Insert(player); err != nil {
			return err
		}
		return nil
	}
}

// 注册多个玩家
func (slf *PlayerEngine) RegistryMulti(unit interfaces.UnitInterface, sort int, playerAccounts ...string) error {
	var players []*models.Player
	var formatAccounts []string
	for _, account := range playerAccounts {
		format := strings.TrimSpace(account)
		players = append(players, models.NewPlayer(
			unit, format, format, format, sort))
		formatAccounts = append(formatAccounts, format)
	}
	var findPlayers []*models.Player
	if count, err := orm.Get().QueryTable("player").
		Filter("is_delete", false).
		Filter("unit_id", unit.GetId()).
		Filter("account__in", formatAccounts).
		All(&findPlayers, "account"); err != nil {
		return err
	} else if count > 0 {
		source := "these players already exist: "
		for _, player := range findPlayers {
			source += player.Account + ", "
		}
		return errors.New(kstr.RemoveLast(kstr.RemoveLast(source)))
	} else {
		rollbackOrm := orm2.NewOrm()
		err := rollbackOrm.Begin()
		if err != nil {
			return err
		}
		_, err = rollbackOrm.InsertMulti(100, players)
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

// 在特定的单位登录玩家，返回token凭证及错误信息
func (slf *PlayerEngine) Login(unit interfaces.UnitInterface, playerAccount, playerPassword string) (tokenString string, e error) {
	var player = new(models.Player)
	if err := orm.Get().QueryTable("player").Filter("is_delete", false).
		Filter("unit_id", unit.GetId()).
		Filter("account", playerAccount).
		Filter("password", playerPassword).Limit(1).One(player); err != nil {
		return "", err
	}

	// 满足登录情况
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	tokenStr, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		logs.Error(err)
		return "", errors.New("签名许可令牌时发生错误")
	}

	// todo:存储用户在线信息
	return tokenStr, nil
}
