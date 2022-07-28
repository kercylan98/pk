package orm

import "github.com/astaxie/beego/orm"

type WheresInterface interface {
	Handle(seter orm.QuerySeter) orm.QuerySeter
}

type Where struct {
	wheres []WheresInterface
}

func (slf *Where) Add(wheres ...WheresInterface) *Where {
	if slf.wheres == nil {
		slf.wheres = make([]WheresInterface, 0)
	}
	slf.wheres = append(slf.wheres, wheres...)
	return slf
}

func (slf *Where) Handle(seter orm.QuerySeter) orm.QuerySeter {
	for _, where := range slf.wheres {
		seter = where.Handle(seter)
	}
	return seter
}