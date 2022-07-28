package interfaces

type UnitTeacherInterface interface {
	// 获取ID
	GetId() string
	// 获取名称
	GetName() string

	// 改变名称
	ChangeName(newName string) error
	// 改变序号
	ChangeSort(newSort int) error

	// 删除教师
	Delete() error
}
