package checks

import (
	"pk/common/conn"
	"unicode/utf8"
)

func CheckUnitName(unitName string) (interface{}, interface{}, bool) {
	if unitName == "" {
		return conn.CODE_ILLEGAL_ARGS, "err:单位名称不允许为空", false
	}
	if utf8.RuneCountInString(unitName) > 30 {
		return conn.CODE_ILLEGAL_ARGS, "err:单位名称不能超过最大30个字", false
	}
	if utf8.RuneCountInString(unitName) > 30 {
		return conn.CODE_ILLEGAL_ARGS, "err:单位名称不能超过最大30个字", false
	}

	return nil, nil, true
}
