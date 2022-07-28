package wheres

import (
	"github.com/astaxie/beego/orm"
	"github.com/kercylan98/dev-kits/utils/k"
	orm2 "pk/common/orm"
)

type gt struct {
	field string
	value interface{}
}

func Gt(field string, comparisonValue interface{}, equal ...bool) orm2.WheresInterface {
	matchEqual := false
	if len(equal) > 0 {
		matchEqual = equal[0]
	}
	return &gt{
		field: k.If(matchEqual, field+"__gte", field+"__gt").(string),
		value: comparisonValue,
	}
}

func (slf *gt) Handle(seter orm.QuerySeter) orm.QuerySeter {
	return seter.Filter(slf.field, slf.value)
}
