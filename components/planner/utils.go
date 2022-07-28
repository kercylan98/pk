package planner

import (
	"github.com/kercylan98/dev-kits/utils/kstr"
	"github.com/kercylan98/dev-kits/utils/kstreams"
	"reflect"
	"strings"
)

// 获取各个课位冲突信息
func GetAllCause(plan *Plan, course *Course) [][]string {
	var causes [][]string
	for i := 0; i < len(plan.Journeys); i++ {
		causes = append(causes, []string{})
		for a := 0; a < len(plan.Journeys[0]); a++ {
			causes[i] = append(causes[i], "")
		}
	}

	for w, week := range plan.Journeys {
		for s, _ := range week {
			cause := ""
			if !checkRoomAllow(plan, course.Room, w, s) {
				cause += "场地冲突,"
			}

			if !checkTeacherAllow(plan, course.Teachers, w, s) {
				cause += "教师冲突</br>"
			} else {
				cause = kstr.RemoveLast(cause) + "</br>"
			}

			allStu := getStudents(plan, w, s)
			m := map[string]*Student{}
			for _, student := range allStu {
				m[student.Number] = student
			}
			allStu = []*Student{}
			kstreams.EachMapSort(m, func(key string, val *Student) {
				allStu = append(allStu, val)
			})

			stuExist := false
			for _, student := range allStu {
				for _, s := range course.Students {
					if student.Number == s.Number {
						cause += student.Name + ","
						stuExist = true
					}
				}
			}

			if stuExist {
				cause = kstr.RemoveLast(cause) + " 学生冲突"
			}
			if strings.HasPrefix(cause, "</br>") {
				cause = cause[5:]
			}
			causes[w][s] = cause
		}
	}
	return causes
}

// 获取特定课程一天内的课时数
func GetWeeklyNumber(plan *Plan, course *Course, week int) int {
	count := 0
	for _, section := range plan.Journeys[week] {
		for _, c := range section {
			if c.Name == course.Name {
				count++
			}
		}
	}
	return count
}

// 特定课位是否有课程冲突
func IsAllow(plan *Plan, course *Course, week, section int) bool {
	for _, c := range plan.Journeys[week][section] {
		if CompareHasConflict(course, c) {
			return false
		}
	}
	return true
}

// 获取特定课程所在课位
func GetSection(plan *Plan, course *Course) [][]int {
	var sections [][]int
	for w, week := range plan.Journeys {
		for s, section := range week {
			if len(sections) < course.WeeklyNumber {
				for _, c := range section {
					if c.Name == course.Name {
						sections = append(sections, []int{w, s})
					}
				}
			} else {
				return sections
			}
		}
	}
	return sections
}

// 得到特定课程在特定课位冲突的课程
func GetConflictCourse(plan *Plan, course *Course, week, section int) []*Course {
	var conflicts []*Course
	for _, c := range plan.Journeys[week][section] {
		if course.Name == c.Name {
			continue
		}
		if CompareHasConflict(course, c) {
			conflicts = append(conflicts, c)
		}
	}
	return conflicts
}

// 比较两个课程是否冲突
func CompareHasConflict(courseA *Course, courseB *Course) bool {
	if courseA.Room == courseB.Room && strings.TrimSpace(courseA.Room) != "" && strings.TrimSpace(courseB.Room) != "" {
		return true
	}
	for _, teacher := range courseA.Teachers {
		for _, t := range courseB.Teachers {
			if teacher == t {
				return true
			}
		}
	}
	for _, student := range courseA.Students {
		for _, s := range courseB.Students {
			if student.Number == s.Number {
				return true
			}
		}
	}
	return false
}

// 得到与特定课程不冲突的所有课程
func GetNoConflictCourse(plan *Plan, course *Course) []*Course {
	collect := plan.Waits
	// 汇总所有课程
	for _, week := range plan.Journeys {
		for _, section := range week {
			for _, course := range section {
				collect = append(collect, course)
			}
		}
	}
	// 匹配
	var match []*Course
	for _, c := range collect {
		if !CompareHasConflict(course, c) {
			match = append(match, c)
		}
	}
	return match
}

// 得到所有可排课位
func GetAllows(plan *Plan, course *Course) [][]int {
	var allAllows [][]int
	for w, week := range plan.Journeys {
		for s, _ := range week {

			// 检查冲突
			if checkStudentsAllow(plan, course.Students, w, s) &&
				checkTeacherAllow(plan, course.Teachers, w, s) &&
				checkRoomAllow(plan, course.Room, w, s) &&
				!isForbidden(plan, course, w, s) {
				allAllows = append(allAllows, []int{w, s})
			}
		}
	}
	return allAllows
}

// 得到所有冲突课位
func GetNotAllow(plan *Plan, course *Course) [][]int {
	var notAllows [][]int
	for w, week := range plan.Journeys {
		for s, _ := range week {

			// 检查冲突
			if checkStudentsAllow(plan, course.Students, w, s) &&
				checkTeacherAllow(plan, course.Teachers, w, s) &&
				checkRoomAllow(plan, course.Room, w, s) &&
				!isForbidden(plan, course, w, s) {

			} else {
				notAllows = append(notAllows, []int{w, s})
			}
		}
	}
	return notAllows
}

// 场地在特定课位是否没有冲突
func checkRoomAllow(plan *Plan, room string, week, section int) bool {
	if strings.TrimSpace(room) == "" {
		return true
	}
	for _, course := range plan.Journeys[week][section] {
		if course.Room == room && strings.TrimSpace(course.Room) != "" {
			return false
		}
	}
	return true
}

// 一组学生在特定课位是否没有冲突
func checkStudentsAllow(plan *Plan, students []*Student, week, section int) bool {
	stus := getStudents(plan, week, section)
	for _, s := range stus {
		for _, student := range students {
			if student.Number == s.Number {
				return false
			}
		}
	}

	return true
}

// 一组老师在特定课位是否没有冲突
func checkTeacherAllow(plan *Plan, teacher []string, week, section int) bool {
	teachers := getTeachers(plan, week, section)
	for _, t := range teachers {
		for _, ts := range teacher {
			if ts == t {
				return false
			}
		}
	}
	return true
}

// 得到特定课位所有老师
func getTeachers(plan *Plan, week, section int) []string {
	teachers := map[string]int{}
	for _, course := range plan.Journeys[week][section] {
		for _, teacher := range course.Teachers {
			teachers[teacher] = 1
		}
	}
	var result []string
	for k, _ := range teachers {
		result = append(result, k)
	}
	return result
}

// 得到特定课位所有学生
func getStudents(plan *Plan, week, section int) []*Student {
	students := map[string]*Student{}
	for _, course := range plan.Journeys[week][section] {
		for _, student := range course.Students {
			students[student.Number] = student
		}
	}
	var result []*Student
	for _, student := range students {
		result = append(result, student)
	}
	return result
}

func GetForbiddenNumber(num interface{}) int {
	return getForbiddenNumber(num)
}

// 获取禁排数值
func getForbiddenNumber(num interface{}) int {
	switch reflect.ValueOf(num).Kind() {
	case reflect.Float64:
		return int(num.(float64))
	case reflect.Int:
		return int(num.(int))
	}
	return num.(int)
}

// 是否是禁排的课位
func isForbidden(plan *Plan, course *Course, week, section int) bool {
	for _, forbidden := range plan.Forbidden {
		var w, s = getForbiddenNumber(forbidden[0]), getForbiddenNumber(forbidden[1])

		if w == week && s == section {
			switch forbidden[2].(string) {
			case FORBIDDEN_DEFAULT:
				return true
			case FORBIDDEN_STAGE, FORBIDDEN_STAGE_CHANGE_WEEK, FORBIDDEN_STAGE_CHANGE_WEEK_NOSAME:
				if course.Stage == forbidden[3].(string) {
					return true
				}
			case FORBIDDEN_TEACHER:
				for _, teacher := range course.Teachers {
					if teacher == forbidden[3].(string) {
						return true
					}
				}
			}

		}
	}
	return false
}