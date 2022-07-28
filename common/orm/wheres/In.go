package wheres

import (
	orm2 "pk/common/orm"
	"github.com/astaxie/beego/orm"
)

type in struct {
	field 		string
	value 		[]interface{}
}

func In(field string, match ...interface{}) orm2.WheresInterface {
	return &in{
		field: field + "__in",
		value: match,
	}
}

func (slf *in) Handle(seter orm.QuerySeter) orm.QuerySeter {
	return seter.Filter(slf.field, slf.value)
}

