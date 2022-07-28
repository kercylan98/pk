package olds

import (
	"fmt"
	"github.com/google/uuid"
	"os"
	"pk/common/basic"
	"pk/common/conn"
	"pk/common/utils"
	"pk/components/planner"
	"strings"
)

type PlanController struct {
	basic.Controller
}

// @Title 新建排课方案
// @Description 创建一个全新的排课方案
// @Param	planName	formData	string	true	"方案名称"
// @Param	week	formData	int	true	"方案上课周次"
// @Param	section	formData	int	true	"方案上课节次"
// @Success 200 无
// @Failure 500 参数不符合规范或方案重复
// @router /new [post]
func (slf *PlanController) NewPlan() {
	slf.Data["json"] = conn.Dispose(func(planName string, week, section int) interface{} {
		if plan, err := planner.NewPlan(planName, week, section); err == nil {
			planner.SetOnlinePlan(slf.Ctx.Request, plan)
			return nil
		}else {
			return fmt.Sprint("err:", err.Error())
		}
	},
		slf.GetString, "planName",
		slf.GetInt, "week",
		slf.GetInt, "section")
	slf.ServeJSON()
}

// @Title 渲染课表
// @Description 渲染课表
// @Success 200 无
// @Failure 500 参数不符合规范或方案重复
// @router /draw [get]
func (slf *PlanController) Draw() {
	slf.Data["json"] = conn.Dispose(func() (interface{}, interface{}) {
		if plan := planner.GetOnlinePlan(slf.Ctx.Request); plan == nil { return "err:请选择排课方案", "code:10000" } else {
			if err := plan.Draw(); err != nil {
				return "err:课表渲染异常：" + err.Error(), "code:500"
			}
			if err := utils.Zip(plan.Name, "assets/" + plan.Name + ".zip"); err != nil {
				return "err:课表渲染异常：" + err.Error(), "code:500"
			}
			slf.Data["planname"] = plan.Name
		}
		return nil, nil
	})
	slf.Ctx.Output.Download("assets/" + slf.Data["planname"].(string) + ".zip", slf.Data["planname"].(string) + "_所有课表.zip")
}

// @Title 数据导入模板下载
// @Description 数据导入模板下载
// @Success 200 无
// @Failure 500 下载失败
// @router /template [get]
func (slf *PlanController) Template() {
	slf.Ctx.Output.Download("assets/temp/导入模板.xlsx")
}

// @Title 切换在线方案
// @Description 切换当前在线方案为特定方案
// @Param	planName	formData	string	true	"需要切换至在线状态的方案名称"
// @Success 200 无
// @Failure 404 找不到特定方案
// @router /switch [post]
func (slf *PlanController) SwitchOnline() {
	slf.Data["json"] = conn.Dispose(func(planName string) (interface{}, interface{}) {
		if plan := planner.GetPlan(strings.TrimSpace(planName)); plan == nil {
			return "err:方案“" + planName + "”不存在", "code:404"
		} else {
			// 保存当前方案
			if nowPlan := planner.GetOnlinePlan(slf.Ctx.Request); nowPlan != nil {
				nowPlan.Save()
			}
			planner.SetOnlinePlan(slf.Ctx.Request, plan)
			return nil, nil
		}
	},
		slf.GetString, "planName")
	slf.ServeJSON()
}

// @Title 添加学生
// @Description 添加学生
// @Success 200 无
// @Failure 10000 没有打开任何排课方案
// @router /add [get]
func (slf *PlanController) Add() {
	slf.Data["json"] = conn.Dispose(func() (interface{}, interface{}) {
		if plan := planner.GetOnlinePlan(slf.Ctx.Request); plan == nil { return "err:请选择排课方案", "code:10000" } else {
			name := slf.GetString("name")
			number := slf.GetString("number")
			course := slf.GetString("course")
			if strings.TrimSpace(name) != "" && strings.TrimSpace(number) != "" && strings.TrimSpace(course) != "" {
				plan.AddStudent(name, number, course)
			}
			return nil, nil
		}
	})
	slf.ServeJSON()
}


// @Title 删除学生
// @Description 删除学生
// @Success 200 无
// @Failure 10000 没有打开任何排课方案
// @router /del [get]
func (slf *PlanController) Del() {
	slf.Data["json"] = conn.Dispose(func() (interface{}, interface{}) {
		if plan := planner.GetOnlinePlan(slf.Ctx.Request); plan == nil { return "err:请选择排课方案", "code:10000" } else {
			number := slf.GetString("number")
			course := slf.GetString("course")
			if strings.TrimSpace(number) != "" && strings.TrimSpace(course) != "" {
				plan.DelStudent(number, course)
			}
			return nil, nil
		}
	})
	slf.ServeJSON()
}

// @Title 自动排课
// @Description 根据当前打开的排课方案的课程、学生等数据进行自动排课
// @Success 200 无
// @Failure 10000 没有打开任何排课方案
// @router /auto [post]
func (slf *PlanController) Auto() {
	slf.Data["json"] = conn.Dispose(func() (interface{}, interface{}) {
		if plan := planner.GetOnlinePlan(slf.Ctx.Request); plan == nil { return "err:请选择排课方案", "code:10000" } else {
			plan.AutoBuild()
			return nil, nil
		}
	})
	slf.ServeJSON()
}

// @Title 自动优化
// @Description 根据当前打开排课方案的各项设定进行冲突优化
// @Success 200 无
// @Failure 10000 没有打开任何排课方案
// @router /optimize [post]
func (slf *PlanController) Optimize() {
	slf.Data["json"] = conn.Dispose(func() (interface{}, interface{}) {
		if plan := planner.GetOnlinePlan(slf.Ctx.Request); plan == nil { return "err:请选择排课方案", "code:10000" } else {
			plan.OptimizeCount = 0
			plan.OptimizeUseless = 0
			plan.Optimize()
			return nil, nil
		}
	})
	slf.ServeJSON()
}

// @Title 排课数据导入
// @Description 通过导入模板进行数据导入
// @Success 200 无
// @Failure 500 服务器内部异常
// @router /import [post]
func (slf *PlanController) Import() {
	slf.Data["json"] = conn.Dispose(func() interface{} {
		f, _, err := slf.GetFile("file")
		if err != nil {
			return fmt.Sprint("err:", err)
		}
		filename := "assets/temp/" + uuid.New().String()
		defer func() {
			f.Close()
			os.Remove(filename)
		}()
		err = slf.SaveToFile("file", filename)
		if err != nil {
			return fmt.Sprint("err:", err)
		}

		if err := planner.LoadData(strings.TrimSpace(slf.GetString("planName")), filename); err != nil {
			return fmt.Sprint("err:", err)
		}else {
			return nil
		}
	})
	slf.ServeJSON()
}

// @Title 移动课位
// @Description 移动课位
// @Success 200 无
// @Failure 500 传入的参数错误或超出范围
// @Failure 10000 没有打开任何排课方案
// @router /section/move [post]
func (slf *PlanController) SectionMove() {
	slf.Data["json"] = conn.Dispose(func(week, section, targetWeek, targetSection int) (interface{}, interface{}) {
		if plan := planner.GetOnlinePlan(slf.Ctx.Request); plan == nil { return "err:请选择排课方案", "code:10000" } else {
			for _, forbidden := range plan.Forbidden {
				if planner.GetForbiddenNumber(forbidden[0]) == targetWeek &&
					planner.GetForbiddenNumber(forbidden[1]) == targetSection{
					return "err:目标课位为禁排课位，需要在Excel调整", "code:500"
				}
			}

			courses := plan.Journeys[week][section]
			targetCourses := plan.Journeys[targetWeek][targetSection]
			plan.Journeys[week][section] = targetCourses
			plan.Journeys[targetWeek][targetSection] = courses
			plan.Save()
			return nil, nil
		}
	},
		slf.GetInt, "week",
		slf.GetInt, "section",
		slf.GetInt, "targetWeek",
		slf.GetInt, "targetSection")
	slf.ServeJSON()
}

// @Title 移动当前打开的方案课程位置
// @Description 将当前打开的方案的特定课程移动到指定的课位
// @Param	courseName	formData	string	true	"需要移动课位的课程名称"
// @Param	week	formData	int	true	"需要移动课位的课程的当前所在周次"
// @Param	section	formData	int	true	"需要移动课位的课程的当前所在节次"
// @Param	targetWeek	formData	int	true	"需要移动课位的课程的目标所在周次"
// @Param	targetSection	formData	int	true	"需要移动课位的课程的目标所在节次"
// @Success 200 无
// @Failure 500 传入的参数错误或超出范围
// @Failure 10000 没有打开任何排课方案
// @router /course/move [post]
func (slf *PlanController) CourseMove() {
	slf.Data["json"] = conn.Dispose(func(courseName string, week, section, targetWeek, targetSection int) (interface{}, interface{}) {
		if plan := planner.GetOnlinePlan(slf.Ctx.Request); plan == nil { return "err:请选择排课方案", "code:10000" } else {
			// 如果目标week和section为-1，且week和section不为-1则放回待排课区域
			if targetWeek == -1 && targetSection == -1 && week != -1 && section != -1 {
				var newSection []*planner.Course
				for _, course := range plan.Journeys[week][section] {
					if course.Name == courseName {
						fmt.Println("方案", plan.Name, "中", courseName, "已调整至待排课区域")
						plan.Waits = append(plan.Waits, course)
					}else {
						newSection = append(newSection, course)
					}
				}
				plan.Journeys[week][section] = newSection
			}else {
				// 如果week或者section为-1，则从待排课位找
				if week == -1 || section == -1 {
					var newWait []*planner.Course
					isHandle := false
					for _, course := range plan.Waits {
						if course.Name == courseName && isHandle == false {
							isHandle = true
							if planner.IsAllow(plan, course, targetWeek, targetSection) {
								plan.Journeys[targetWeek][targetSection] = append(plan.Journeys[targetWeek][targetSection], course)
							}else {
								return "err:由于存在冲突，无法将课程调整到该位置", "code:500"
							}
						}else {
							newWait = append(newWait, course)
						}
					}
					plan.Waits = newWait
				}else {
					var newSection []*planner.Course
					for _, course := range plan.Journeys[week][section] {
						if course.Name == courseName {
							if planner.IsAllow(plan, course, targetWeek, targetSection) {
								fmt.Println("方案", plan.Name, "中", courseName, "已调整至周", targetWeek, "第", targetSection, "节")
								plan.Journeys[targetWeek][targetSection] = append(plan.Journeys[targetWeek][targetSection], course)
							}else {
								return "err:由于存在冲突，无法将课程调整到该位置", "code:500"
							}
						}else {
							newSection = append(newSection, course)
						}
					}
					plan.Journeys[week][section] = newSection
				}
			}

			plan.Save()
			return nil, nil
		}
	},
		slf.GetString, "courseName",
		slf.GetInt, "week",
		slf.GetInt, "section",
		slf.GetInt, "targetWeek",
		slf.GetInt, "targetSection")
	slf.ServeJSON()
}

// @Title 获取特定课程不冲突的课位
// @Description 在当前打开的排课方案中寻找特定课程不冲突课位并返回
// @Param	courseName	formData	string	true	"需要寻找的课程名称"
// @Success 200 {[][]int} []int{week,section}
// @Failure 500 传入的参数错误或超出范围
// @Failure 10000 没有打开任何排课方案
// @router /course/allows [get]
func (slf *PlanController) Allows() {
	slf.Data["json"] = conn.Dispose(func(courseName string) (interface{}, interface{}) {
		courseName = strings.TrimSpace(courseName)
		if plan := planner.GetOnlinePlan(slf.Ctx.Request); plan == nil { return "err:请选择排课方案", "code:10000" } else {
			for _, course := range plan.Waits {
				if course.Name == courseName {
					course.Cause = planner.GetAllCause(plan, course)
					return planner.GetAllows(plan, course), course
				}
			}
		}
		return fmt.Sprint("err:不存在的待排课程 ", courseName), "code:404"
	},
		slf.GetString, "courseName")
	slf.ServeJSON()
}

// @Title 获取特定课程存在冲突的课位
// @Description 在当前打开的排课方案中寻找特定课程存在冲突的课位并返回
// @Param	courseName	formData	string	true	"需要寻找的课程名称"
// @Param	week	formData	int	true	"当前所在周次"
// @Param	section	formData	int	true	"当前所在节次"
// @Success 200 {[][]int} []int{week,section}
// @Failure 500 传入的参数错误或超出范围
// @Failure 10000 没有打开任何排课方案
// @router /course/unallowable [get]
func (slf *PlanController) Unallowable() {
	slf.Data["json"] = conn.Dispose(func(courseName string, week, section int) (interface{}, interface{}) {
		if plan := planner.GetOnlinePlan(slf.Ctx.Request); plan == nil { return "err:请选择排课方案", "code:10000" } else {
			for _, course := range plan.Journeys[week][section] {
				if course.Name == courseName {
					course.Cause = planner.GetAllCause(plan, course)
					return planner.GetNotAllow(plan, course), course
				}
			}

		}
		return fmt.Sprint("err:未找到课程 ", courseName), "code:404"
	},
		slf.GetString, "courseName",
		slf.GetInt, "week",
		slf.GetInt, "section")
	slf.ServeJSON()
}
