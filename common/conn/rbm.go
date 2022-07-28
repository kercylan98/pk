package conn

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"reflect"
	"strconv"
	"strings"
)

// 响应文内容数据结构模型
type ResponseBodyModel struct {
	Code  				int           				// 请求结果状态码
	Data				[]interface{} 				// 请求结果数据,默认应该使用0号数据
	Error 				string						// 发生的异常
}

// 处理内容
func Dispose(disposeFunc interface{}, checkFuncAndKey ...interface{}) *ResponseBodyModel {
	rbm := &ResponseBodyModel{
		Code:   0,
		Data: 	[]interface{}{},
	}
	// 解析checkFunc，将函数和参数key区分开来
	// 不管函数和key的顺序如何，第一个函数应该取到的始终是第一个key
	var allFunc, allKey = make([]reflect.Value, 0), make([]reflect.Value, 0)
	for _, check := range checkFuncAndKey {
		v := reflect.ValueOf(check)
		if v.Kind() == reflect.Func {
			allFunc = append(allFunc, v)
		}else {
			allKey = append(allKey, v)
		}
	}
	// 确保函数和键长度一致
	if len(allFunc) != len(allKey) {
		err := errors.New("the validation method must receive the same number of functions as the number of key names")
		logs.Error(err)
		rbm.Code = -1
		rbm.Error = "系统内部发生未预判到的错误"
		return rbm
	}else {
		// 处理函数入参
		var disposeArgs []reflect.Value
		// 解析参数
		for i := 0; i < len(allFunc); i++ {
			f, k := allFunc[i], allKey[i]
			// 调用校验函数解析返回值
			// 如果不是错误则添加到调用处理函数的参数表
			for _, resultValue := range f.Call([]reflect.Value{k}) {
				format := fmt.Sprint(resultValue)
				if strings.HasPrefix(format, "strconv.Atoi") {
					rbm.Error = errors.New(fmt.Sprint("invalid parameter \"", k, "\"")).Error()
					rbm.Code = -1
					return rbm
				}else {
					if format != "<nil>" {
						disposeArgs = append(disposeArgs, resultValue)
					}
				}
			}
		}
		// 调用处理函数，允许处理函数返回值来影响结果
		// 当返回值类型为 "code:500"则表示设置返回值状态码，其他结果均按顺序添加至Data中
		disposeFuncValue := reflect.ValueOf(disposeFunc)
		if len(disposeArgs) != disposeFuncValue.Type().NumIn() {
			logs.Error(errors.New("Handler crashed with error reflect: Call with too few input arguments"))
			rbm.Code = -1
			rbm.Error = "系统内部发生未预判到的错误"
			return rbm
		}
		result := disposeFuncValue.Call(disposeArgs)
		for _, value := range result {
			format := strings.TrimSpace(fmt.Sprint(value))
			if format == "<nil>" {
				continue
			}
			if strings.HasPrefix(format, "code:") {
				if code, err := strconv.Atoi(strings.ReplaceAll(format, "code:", "")); err != nil{
					logs.Error(err)
					rbm.Code = -1
					rbm.Error = "系统内部发生未预判到的错误"
					return rbm
				}else {
					rbm.Code = code
					continue
				}
			}
			if strings.HasPrefix(format, "err:") {
				err := errors.New(strings.TrimSpace(format[4:]))
				logs.Error(err)
				if rbm.Code == 0 {
					rbm.Code = -1
				}
				if format == CODE_SYSTEM_ERR {
					logs.Error(err)
					rbm.Error = errors.New("系统内部发生未预判到的错误").Error()
				}else {
					rbm.Error = err.Error()
				}
				continue
			}
			rbm.Data = append(rbm.Data, value.Interface())
		}

		if rbm.Error == "" && rbm.Code == 0 {
			rbm.Error = "未发生错误"
		}
		return rbm
	}
}

