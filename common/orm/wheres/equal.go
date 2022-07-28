package wheres

import (
	"github.com/astaxie/beego/orm"
	"github.com/kercylan98/dev-kits/utils/k"
	orm2 "pk/common/orm"
)

type equal struct {
	field string
	value interface{}
}

func Equal(field string, value interface{}, caseSensitive ...bool) orm2.WheresInterface {
	var sensitive = true
	if len(caseSensitive) > 0 {
		sensitive = caseSensitive[0]
	}
	return &equal{
		field: k.If(sensitive, field, field+"__iexact").(string),
		value: value,
	}
}

func (slf *equal) Handle(seter orm.QuerySeter) orm.QuerySeter {
	return seter.Filter(slf.field, slf.value)
}
