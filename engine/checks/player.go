package checks

import (
	"pk/common/conn"
	"strings"
	"unicode/utf8"
)

func CheckPlayerAccount(account string) (interface{}, interface{}, bool) {
	account = strings.TrimSpace(account)
	if account == "" {
		return conn.CODE_ILLEGAL_ARGS, "err:登录帐号不允许为空", false
	}
	if utf8.RuneCountInString(account) > 20 {
		return conn.CODE_ILLEGAL_ARGS, "err:登录帐号长度不能超过20个字符", false
	}
	if utf8.RuneCountInString(account) < 3 {
		return conn.CODE_ILLEGAL_ARGS, "err:登录帐号长度不能小于3个字符", false
	}

	return nil, nil, true
}

func CheckPlayerPassword(password string) (interface{}, interface{}, bool) {
	if password == "" {
		return conn.CODE_ILLEGAL_ARGS, "err:登录密码不允许为空", false
	}
	if utf8.RuneCountInString(password) > 20 {
		return conn.CODE_ILLEGAL_ARGS, "err:登录帐号长度不能超过30个字符", false
	}
	if utf8.RuneCountInString(password) < 6 {
		return conn.CODE_ILLEGAL_ARGS, "err:登录帐号长度不能小于6个字符", false
	}

	return nil, nil, true
}


func CheckPlayerName(name string) (interface{}, interface{}, bool) {
	if name == "" {
		return conn.CODE_ILLEGAL_ARGS, "err:账户名称不允许为空", false
	}
	if utf8.RuneCountInString(name) > 15 {
		return conn.CODE_ILLEGAL_ARGS, "err:登录帐号长度不能超过15个字符", false
	}

	return nil, nil, true
}