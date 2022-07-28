package wheres

import (
	"github.com/astaxie/beego/orm"
	"github.com/kercylan98/dev-kits/utils/k"
	orm2 "pk/common/orm"
)

type lt struct {
	field string
	value interface{}
}

func Lt(field string, comparisonValue interface{}, equal ...bool) orm2.WheresInterface {
	matchEqual := false
	if len(equal) > 0 {
		matchEqual = equal[0]
	}
	return &lt{
		field: k.If(matchEqual, field+"__lte", field+"__lt").(string),
		value: comparisonValue,
	}
}

func (slf *lt) Handle(seter orm.QuerySeter) orm.QuerySeter {
	return seter.Filter(slf.field, slf.value)
}
