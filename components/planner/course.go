package planner

// 课程数据模型结构
type Course struct {
	Number 					string   // 课程编号
	Name 					string     // 课程名称
	Teachers				[]string   // 课程任教
	WeeklyNumber 			int        // 周课时数
	Students 				[]*Student // 课程学生
	Stage 					string    // 课程阶段
	Room 					string     // 课程教室
	TeacherStr 				string   // 教师字符串

	Cause 					[][]string				// 各课位冲突原因
	IsRemoveWait			bool					// 是否需要从待排课区域移除
}
