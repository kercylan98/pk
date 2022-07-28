package orm

import (
	"errors"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"sync"
)

// 全局Ormer通用实例
// 进行事物操作时应另外生成
var globalOrmerInstanceOnce sync.Once
var globalOrmerInstance orm.Ormer

var (
	ErrTxHasBegan    = errors.New("<Ormer.Begin> transaction already begin")
	ErrTxDone        = errors.New("<Ormer.Commit/Rollback> transaction not begin")
	ErrMultiRows     = errors.New("<QuerySeter> return multi rows")
	ErrNoRows        = errors.New("<QuerySeter> no row found")
	ErrStmtClosed    = errors.New("<QuerySeter> stmt already closed")
	ErrArgs          = errors.New("<Ormer> args error may be empty")
	ErrNotImplement  = errors.New("have not implement")
)

// 获取全局Ormer连接
func Get() orm.Ormer {
	globalOrmerInstanceOnce.Do(func() {
		globalOrmerInstance = orm.NewOrm()
	})
	return globalOrmerInstance
}

// 创建临时的Ormer连接
func New(aliasName ...string) orm.Ormer {
	o :=  orm.NewOrm()
	if len(aliasName) > 0{
		o.Using(aliasName[0])
	}
	return o
}

// 错误比较
func EqErr(a, b error) bool {
	return strings.TrimSpace(a.Error()) == strings.TrimSpace(b.Error())
}