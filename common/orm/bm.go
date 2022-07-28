package orm


// 所有数据模型的基类
type Model struct {
	IsDelete		bool					`orm:"default(false)" description:"是否已删除"`
}

// 如果数组存在内容则返回第一个元素，否则返回默认值
func HasGetOneOrElseDefault(slice []int, defaultValue int) int {
	if len(slice) > 0 {
		return slice[0]
	}
	return defaultValue
}

func (slf *Model) Del() error {
	slf.IsDelete = true
	_, err := Get().Update(slf, "IsDelete")
	return err
}
