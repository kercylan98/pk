package wheres

import (
	"github.com/astaxie/beego/orm"
	"github.com/kercylan98/dev-kits/utils/k"
	orm2 "pk/common/orm"
)

type contains struct {
	field string
	value interface{}
}

func Contains(field string, value interface{}, caseSensitive ...bool) orm2.WheresInterface {
	var sensitive = true
	if len(caseSensitive) > 0 {
		sensitive = caseSensitive[0]
	}
	return &contains{
		field: k.If(sensitive, field+"__contains", field+"__icontains").(string),
		value: value,
	}
}

func (slf *contains) Handle(seter orm.QuerySeter) orm.QuerySeter {
	return seter.Filter(slf.field, slf.value)
}
