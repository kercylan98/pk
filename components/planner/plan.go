package planner

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/kercylan98/dev-kits/utils/krand"
	"github.com/tealeg/xlsx"
	"os"
	"strings"
	"sync"
)

// 排课计划数据模型
type Plan struct {
	sync.RWMutex `json:"-"`
	Name         string        // 排课方案名称
	Waits        []*Course     // 待安排的课程信息
	Journeys     [][][]*Course // 课程详情安排信息
	Rooms        []string      // 所有可用教室

	StudentHasCourse map[string][]*Course // 学生拥有的课程（Key：学号）
	OptimizeCount    int                  // 优化次数
	MaxOptimizeCount int                  // 最大优化次数
	CourseCount      int                  // 总课时数
	NowCourseTotal   int                  // 当前总课时数

	Forbidden         [][]interface{} // 禁排课位[]
	IsRunning         bool            // 耗时操作中
	TryContinuous     bool            // 尝试课程尽量连排
	MustContinuousTwo []string        // 必须两节连排的科目
	WeekNumRangMax    int             // 尽量保证每天同一课程不超过多少节，0不限

	Stage []string // 所有年级

	OptimizeUseless int // 优化课表无用功次数
}

// 测试功能
func (slf *Plan) TestFunc() {
	slf.Forbidden = [][]interface{}{}
	slf.TryContinuous = true
	slf.WeekNumRangMax = 2
	// 上午

	slf.AddForbidden(1, 5, FORBIDDEN_DEFAULT)
	slf.AddForbidden(1, 8, FORBIDDEN_DEFAULT)
	slf.AddForbidden(1, 10, FORBIDDEN_DEFAULT)
	slf.AddForbidden(1, 11, FORBIDDEN_DEFAULT)

	slf.AddForbidden(2, 3, FORBIDDEN_DEFAULT)
	slf.AddForbidden(2, 4, FORBIDDEN_DEFAULT)
	slf.AddForbidden(2, 5, FORBIDDEN_DEFAULT)
	slf.AddForbidden(2, 8, FORBIDDEN_DEFAULT)
	slf.AddForbidden(2, 9, FORBIDDEN_DEFAULT)
	slf.AddForbidden(2, 10, FORBIDDEN_DEFAULT)
	slf.AddForbidden(2, 11, FORBIDDEN_DEFAULT)

	slf.AddForbidden(3, 3, FORBIDDEN_DEFAULT)
	slf.AddForbidden(3, 4, FORBIDDEN_DEFAULT)
	slf.AddForbidden(3, 5, FORBIDDEN_DEFAULT)
	slf.AddForbidden(3, 9, FORBIDDEN_DEFAULT)
	slf.AddForbidden(3, 10, FORBIDDEN_DEFAULT)
	slf.AddForbidden(3, 11, FORBIDDEN_DEFAULT)

	slf.AddForbidden(4, 4, FORBIDDEN_DEFAULT)
	slf.AddForbidden(4, 5, FORBIDDEN_DEFAULT)
	slf.AddForbidden(4, 10, FORBIDDEN_DEFAULT)
	slf.AddForbidden(4, 11, FORBIDDEN_DEFAULT)

	slf.AddForbidden(5, 5, FORBIDDEN_DEFAULT)
	slf.AddForbidden(5, 7, FORBIDDEN_DEFAULT)
	slf.AddForbidden(5, 9, FORBIDDEN_DEFAULT)
	slf.AddForbidden(5, 10, FORBIDDEN_DEFAULT)
	slf.AddForbidden(5, 11, FORBIDDEN_DEFAULT)

	//for _, t := range []string{"Scarlett He","Catherine Yang","new Chinese teacher"," New Math teacher","newEnglish teacher","Ally Zhu"} {
	// 数学教研活动
	//slf.AddForbidden(1, 1, FORBIDDEN_TEACHER, t)
	//slf.AddForbidden(2, 1, FORBIDDEN_TEACHER, t)
	//slf.AddForbidden(3, 1, FORBIDDEN_TEACHER, t)
	//slf.AddForbidden(4, 1, FORBIDDEN_TEACHER, t)
	//slf.AddForbidden(5, 1, FORBIDDEN_TEACHER, t)
	//slf.AddForbidden(3, 11, FORBIDDEN_TEACHER, t)
	//}

	for _, t := range []string{"康彬", "沈波", "何思佳", "祝琴琴", "罗倩", "宗慧", "刘娜", "郁龙", "新数学老师", "杨子青", "新语文老师", "新中文老师"} {
		// 语文教研活动
		slf.AddForbidden(2, 1, FORBIDDEN_TEACHER, t)
		slf.AddForbidden(3, 1, FORBIDDEN_TEACHER, t)
		slf.AddForbidden(4, 1, FORBIDDEN_TEACHER, t)
		slf.AddForbidden(5, 1, FORBIDDEN_TEACHER, t)
		//slf.AddForbidden(4, 10, FORBIDDEN_TEACHER, t)
		//slf.AddForbidden(4, 11, FORBIDDEN_TEACHER, t)
	}

	// 尽量连排，同时尽量保证每个课程每天不超过2节
	slf.TryContinuous = true
	slf.WeekNumRangMax = 2

	// DP2 VA: HL 必须2节连排
	slf.MustContinuousTwo = append(slf.MustContinuousTwo, "G10RoseA Optional class Pre IB VA")
	slf.MustContinuousTwo = append(slf.MustContinuousTwo, "G10RoseB Optional class Pre A level Art")
	slf.MustContinuousTwo = append(slf.MustContinuousTwo, "G10Ginkgo Optional class Pre IB Music")
	slf.MustContinuousTwo = append(slf.MustContinuousTwo, "G9Beech Compulsory Art & music（choose 1 out of 2)")
	slf.MustContinuousTwo = append(slf.MustContinuousTwo, "G9Rose Compulsory Art & Design  & music （choose 1 out of 2)")
	//
	//
	//		// 所有课程必须2节连排
	//		slf.MustContinuousTwo = append(slf.MustContinuousTwo, strings.Split(`DP2 physics HL+SL+A-level
	//DP2 physics HL+A-level
	//G12 physics A-level
	//DP2 chemistry HL+SL
	//DP2 chemistry HL
	//DP2 Biology HL&SL
	//DP2 Biology HL
	//DP2 ESS
	//DP2 Sports sceince SL
	//G12 Alevel biology
	//DP2 TOK in Chinese
	//DP2 TOK in English
	//G12 Alevel Critical thinking and writing (CTW)
	//DP2 English B: HL
	//DP2 English B: SL
	//G12 English (AS Level 9093)
	//Combined A-level English
	//DP2 Eco: HL&SL
	//DP2 Eco: HL
	//DP2 BM: HL & SL
	//DP2 BM: HL
	//DP2 Psychology: HL &SL
	//DP2 Psychology: HL
	//DP2 Alevel Economics
	//DP2 Chinese A: literature HL & SL
	//DP2 Chinese A: literature HL
	//DP2 Chinese A:language &literature SL & HL
	//DP2 Chinese A:language &literature HL
	//DP2 Math AA: HL
	//DP2 Math AA: SL
	//G12 A-level Math (Pure)
	//G12 A-level Math (Prob&Stats)
	//DP2 VA: HL
	//DP1 physics HL + SL
	//DP1 physics HL
	//DP1 chemistry HL/SL + A level
	//DP1 chemistry HL + A level
	//G11 chemistry A level
	//DP1 Biology HL/SL
	//G11 Biology A level
	//DP1 Sports sceince HL+SL
	//DP1 Sports sceince HL
	//DP1 ESS
	//G11 Alevel Critical thinking and writing (CTW) &IPQ
	//DP1 English B: HL
	//DP1 English B: SL
	//DP1 English B: SL Rebecca
	//DP1 English B: SL Tina
	//DP1 Eco: HL&SL
	//DP1 Eco: HL
	//DP1 BM: HL & SL
	//DP1 BM: HL
	//G11 A level BM
	//DP1 Psychology: HL &SL
	//DP1 Psychology: HL
	//G11 A level Economics
	//G11 A level Psychology
	//DP1 Chinese A: literature HL & SL
	//DP1 Chinese A: literature HL
	//DP1 Chinese A:language &literature SL A&B
	//DP1 Chinese A:language &literature SL X&Y
	//DP1 Chinese A:language &literature HL
	//DP1 Math AA: HL
	//DP1 Math AA: SL
	//G11 Alevel math
	//DP1 VA: HL +SL
	//DP1 VA: HL
	//DP1 Music: HL +SL
	//DP1 Music: HL`, "\n")...)
	//

	//
	//slf.AddForbidden(1, 1, FORBIDDEN_DEFAULT)
	//slf.AddForbidden(5, 11, FORBIDDEN_DEFAULT)
	//
	//
	// 禁牌课位，周五7、8、9、10节，周二10节、周四10节
	slf.AddForbidden(1, 1, FORBIDDEN_DEFAULT)
	slf.AddForbidden(2, 10, FORBIDDEN_DEFAULT)
	slf.AddForbidden(4, 10, FORBIDDEN_DEFAULT)
	slf.AddForbidden(5, 8, FORBIDDEN_DEFAULT)
	slf.AddForbidden(5, 9, FORBIDDEN_DEFAULT)
	slf.AddForbidden(5, 10, FORBIDDEN_DEFAULT)
	//slf.AddForbidden(2, 10, FORBIDDEN_DEFAULT)
	//slf.AddForbidden(4, 10, FORBIDDEN_DEFAULT)
	//
	//// 全年级School counselling课
	//slf.AddForbidden(1, 8, FORBIDDEN_DEFAULT)
	//
	// 给TOK预留的课位，禁排高二周一9、10节
	slf.AddForbidden(2, 8, FORBIDDEN_STAGE, "G10A")
	slf.AddForbidden(2, 9, FORBIDDEN_STAGE, "G10A")
	slf.AddForbidden(4, 8, FORBIDDEN_STAGE, "G10A")
	slf.AddForbidden(4, 9, FORBIDDEN_STAGE, "G10A")
	slf.AddForbidden(2, 8, FORBIDDEN_STAGE, "G10B")
	slf.AddForbidden(2, 9, FORBIDDEN_STAGE, "G10B")
	slf.AddForbidden(4, 8, FORBIDDEN_STAGE, "G10B")
	slf.AddForbidden(4, 9, FORBIDDEN_STAGE, "G10B")
	slf.AddForbidden(2, 8, FORBIDDEN_STAGE, "G9B")
	slf.AddForbidden(2, 9, FORBIDDEN_STAGE, "G9B")
	slf.AddForbidden(4, 8, FORBIDDEN_STAGE, "G9B")
	slf.AddForbidden(4, 9, FORBIDDEN_STAGE, "G9B")
	slf.AddForbidden(2, 8, FORBIDDEN_STAGE, "G9R")
	slf.AddForbidden(2, 9, FORBIDDEN_STAGE, "G9R")
	slf.AddForbidden(4, 8, FORBIDDEN_STAGE, "G9R")
	slf.AddForbidden(4, 9, FORBIDDEN_STAGE, "G9R")
	//slf.AddForbidden(1, 10, FORBIDDEN_STAGE, "DP2")
	//slf.AddForbidden(1, 8, FORBIDDEN_STAGE, "DP2")

	//teas := map[string]string{}

	//teas[`王棋`]=`1,5|1,8|2,5|3,3|3,4|3,8|3,10|4,1|4,2|4,3|4,5|4,8|5,1|`
	//teas[`EMMETT MARTIN O'BRIEN`]=`1,3|1,4|1,6|1,7|2,3|2,4|3,4|3,5|4,2|4,3|4,7|4,8|5,3|5,4|`
	//teas[`高艳芳`]=`1,3|1,4|1,7|2,5|2,8|2,9|`
	//teas[`刘传云`]=`1,3|1,6|2,3|2,4|2,5|3,1|3,8|3,9|4,3|4,9|5,1|5,2|`
	//teas[`Rebecca Hebert`]=`1,3|1,4|1,8|2,1|2,8|2,9|3,3|3,8|3,9|4,5|4,8|4,9|5,4|5,5|`
	//teas[`DAVID BRUCE WINDSOR BROWN`]=`1,3|1,4|1,8|2,1|2,8|2,9|3,3|3,8|3,9|4,5|4,8|4,9|5,4|5,5|`
	//teas[`刘娜`]=`1,3|1,4|1,8|2,1|2,8|2,9|3,3|3,8|3,9|4,5|4,8|4,9|5,4|5,5|`
	//teas[`HUGH MARCUS BOND`]=`1,3|1,4|1,8|2,1|2,8|2,9|3,3|3,8|3,9|4,5|4,8|4,9|5,4|5,5|`
	//teas[`WADE MORGAN WERNER`]=`1,3|1,4|1,8|2,1|2,8|2,9|3,3|3,8|3,9|4,5|4,8|4,9|5,4|5,5|`
	//teas[`刘益君`]=`1,6|2,5|3,2|4,4|`
	//teas[`外教X`]=`1,4|1,5|1,8|2,2|2,3|2,8|3,1|3,2|3,5|3,6|3,7|3,8|3,10|4,2|4,3|4,7|4,8|5,1|5,2|5,3|5,6|5,7|`
	//teas[`罗倩`]=`2,4|3,10|4,1|5,1|`
	//teas[`mocha(师悦-司机)`]=`2,1|3,2|3,3|4,1|4,2|4,7|5,2|5,3|5,6|5,7|`
	//teas[`陶真`]=`3,2|4,1|4,2|5,3|`
	//teas[`吴磊`]=`1,7|1,8|2,2|2,3|3,1|3,2|3,6|4,1|4,2|4,3|4,7|5,2|5,3|5,7|`
	//teas[`mocha(教师)`]=`3,5|4,4|5,5|`
	//teas[`mocha(师悦)`]=`1,6|2,3|2,8|3,3|3,4|3,5|4,2|4,6|5,6|5,7|`
	//teas[`旷涛群`]=`2,1|2,2|2,3|3,4|3,6|3,7|4,3|4,4|4,8|5,3|5,4|5,5|`
	//teas[`张丽苹`]=`1,7|2,2|3,5|4,4|4,9|5,6|`
	//teas[`朱雪华`]=`1,6|2,2|3,1|3,7|4,4|4,6|`
	//teas[`蔡洋`]=`1,6|1,7|2,4|2,5|3,6|3,7|3,9|3,10|4,6|4,7|5,1|5,2|5,3|5,4|`
	//teas[`Jessie`]=`2,4|2,5|2,9|4,4|4,7|5,1|`
	//teas[`钱亚雯`]=`3,1|3,7|3,9|5,2|5,6|5,7|`
	//teas[`Jessica Hart`]=`1,5|2,1|3,1|3,10|4,5|5,6|`
	//teas[`沈波`]=`2,2|3,5|4,9|`
	//teas[`AHMAD WALI`]=`1,5|2,4|2,9|3,4|3,6|3,7|3,10|4,4|4,5|4,7|5,4|5,5|5,7|`
	//teas[`康彬`]=`2,1|3,2|3,10|4,1|4,6|5,1|5,2|`
	//teas[`NIKOLAOS BITOS`]=`2,4|3,1|3,4|3,6|3,7|4,1|4,6|`
	//teas[`朱雅琪`]=`1,7|1,8|2,2|2,3|3,4|3,6|4,3|4,6|5,6|5,7|`
	//teas[`杨静`]=`1,5|1,6|1,7|`
	//
	//for name, str := range teas {
	//	for _, s := range strings.Split(str, "|") {
	//		if strings.TrimSpace(s) != "" {
	//			ws := strings.Split(s, ",")
	//			w, errw := strconv.Atoi(ws[0])
	//			if errw != nil {
	//				panic(errw)
	//			}
	//			s, errs := strconv.Atoi(ws[1])
	//			if errs != nil {
	//				panic(errs)
	//			}
	//			slf.AddForbidden(w, s, FORBIDDEN_TEACHER, name)
	//		}
	//	}
	//
	//
	//}

}

// 添加学生到特定课程
func (slf *Plan) AddStudent(name, number, course string) {
	for _, journey := range slf.Journeys {
		for _, courses := range journey {
			for _, c := range courses {
				if c.Name == course {
					exist := false
					for _, student := range c.Students {
						if strings.TrimSpace(student.Number) == strings.TrimSpace(number) {
							exist = true
							break
						}
					}
					if !exist {
						c.Students = append(c.Students, &Student{
							Number: strings.TrimSpace(number),
							Name:   strings.TrimSpace(name),
						})

						// 添加学生对应课程缓存数据
						slf.StudentHasCourse[number] = append(slf.StudentHasCourse[number], c)
					}
				}
			}
		}
	}
	slf.Save()
}

// 删除特定课程的学生
func (slf *Plan) DelStudent(number, course string) {
	for _, journey := range slf.Journeys {
		for _, courses := range journey {
			for _, c := range courses {
				if c.Name == course {
					exist := false
					index := 0
					for i, student := range c.Students {
						if strings.TrimSpace(student.Number) == strings.TrimSpace(number) {
							exist = true
							index = i
							break
						}
					}
					if exist {
						newStudents := make([]*Student, 0)
						for i, student := range c.Students {
							if i != index {
								newStudents = append(newStudents, student)
							}
						}
						c.Students = newStudents

						// 添加学生对应课程缓存数据
						newCache := make([]*Course, 0)
						for _, c2 := range slf.StudentHasCourse[number] {
							if c2.Name != c.Name {
								newCache = append(newCache, c2)
							}
						}
						slf.StudentHasCourse[number] = newCache
					}
				}
			}
		}
	}
	slf.Save()
}

// 创建排课方案
func NewPlan(name string, week, section int) (*Plan, error) {
	if _, ok := this.plans[name]; ok {
		return nil, errors.New("plan existed")
	}

	if week <= 0 || section <= 0 {
		return nil, errors.New("week and section must > 0")
	}
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("plan name must is not null")
	}

	plan := &Plan{
		Name:      name,
		Waits:     []*Course{},
		Journeys:  [][][]*Course{},
		Rooms:     []string{},
		Forbidden: [][]interface{}{},

		MaxOptimizeCount: 50,
		StudentHasCourse: map[string][]*Course{},

		Stage: []string{},
	}

	for i := 0; i < week; i++ {
		plan.Journeys = append(plan.Journeys, [][]*Course{})
		for a := 0; a < section; a++ {
			plan.Journeys[i] = append(plan.Journeys[i], []*Course{})
		}
	}

	plan.TestFunc()
	AddPlan(plan)
	return plan, nil
}

// 自动优化
func (slf *Plan) Optimize() {
	waitNumber := len(slf.Waits)

	slf.IsRunning = true
	slf.OptimizeCount++
	fmt.Println("正在进行课表第", slf.OptimizeCount, "次优化。当前剩余冲突数：", len(slf.Waits))

	temp := &Plan{
		Waits:     slf.Waits,
		Journeys:  slf.Journeys,
		Rooms:     slf.Rooms,
		Forbidden: slf.Forbidden,
	}

	// 优先级重排序
	var frist []*Course
	var other []*Course
	for _, wait := range slf.Waits {
		isFrist := false
		for _, s := range slf.MustContinuousTwo {
			if wait.Name == s {
				isFrist = true
				break
			}
		}
		if isFrist {
			frist = append(frist, wait)
		} else {
			other = append(other, wait)
		}
	}
	slf.Waits = append(frist, other...)

	for _, course := range slf.Waits {
		for w, week := range slf.Journeys {
			for s, section := range week {
				for _, usedCourse := range section {

					// 忽略课程
					eq := false
					for _, s := range slf.MustContinuousTwo {
						if usedCourse.Name == s {
							eq = true
							break
						}
					}
					if eq {
						continue
					}

					// 如果这个课程可以放到其他位置，则放过去
					if allows := GetAllows(slf, usedCourse); len(allows) > 0 {
						target := allows[krand.Int(0, len(allows))]

						// 尽量保证每天特定课程不超过多少节
						if slf.WeekNumRangMax > 0 {
							var newAllows [][]int
							for _, ws := range allows {
								if GetWeeklyNumber(slf, course, ws[0]) <= slf.WeekNumRangMax-1 {
									newAllows = append(newAllows, ws)
								}
							}
							if len(newAllows) > 0 {
								allows = newAllows
							}
						}

						//必须两节连排处理
						if len(slf.MustContinuousTwo) > 0 {
							for _, name := range slf.MustContinuousTwo {
								if name == course.Name {
									placed := GetSection(slf, course)
									// 如果允许课程存在于已放置上或下课位  待定：且上上下下课位不存在则使用该课位
									for _, allow := range allows {
										for _, place := range placed {
											if allow[0] == place[0] {
												if allow[1] == place[1]-1 || allow[1] == place[1]+1 {
													target = allow
												}

											}
										}
									}
									break
								}
							}
						}

						// 满足尽量连排，查找这个课程所在课位
						if slf.TryContinuous {
							placed := GetSection(slf, course)
							// 如果允许课程存在于已放置上或下课位  待定：且上上下下课位不存在则使用该课位
							for _, allow := range allows {
								for _, place := range placed {
									if allow[0] == place[0] {
										if allow[1] == place[1]-1 || allow[1] == place[1]+1 {
											target = allow
										}

									}
								}
							}
						}

						// 把课程拿出来
						var newSection []*Course
						for _, tempCourse := range temp.Journeys[w][s] {
							if tempCourse.Name != usedCourse.Name {
								newSection = append(newSection, tempCourse)
							}
						}
						temp.Journeys[w][s] = newSection
						temp.Journeys[target[0]][target[1]] = append(temp.Journeys[target[0]][target[1]], usedCourse)
					}

					// 检查这个待排课程是否可以排课
					if allows := GetAllows(temp, course); len(allows) > 0 {
						target := allows[krand.Int(0, len(allows))]
						temp.Journeys[target[0]][target[1]] = append(temp.Journeys[target[0]][target[1]], course)

						// 从冲突课程内移除
						var newWaits []*Course
						var isHandle = false
						for _, waitCourse := range temp.Waits {
							if waitCourse.Name != course.Name {
								newWaits = append(newWaits, waitCourse)
							} else {
								if !isHandle {
									isHandle = true
								} else {
									newWaits = append(newWaits, waitCourse)
								}
							}
						}
						temp.Waits = newWaits
						goto end
					}
				}
			}
		}
	}

end:
	{
		slf.Waits = temp.Waits
		slf.Journeys = temp.Journeys
		slf.Save()
	}

	if len(slf.Waits) == 0 {
		fmt.Println("课表优化结束，当前剩余冲突：", len(slf.Waits))
		slf.Save()
		slf.IsRunning = false
		return
	}

	if len(slf.Waits) > 0 && slf.OptimizeCount < slf.MaxOptimizeCount {
		if waitNumber == len(slf.Waits) {
			slf.OptimizeUseless++
			// 当无用功次数到达2次的时候就重新抉择禁排课位尝试
			// 当接近优化结尾的时候不进行该行为，确保结果为当前找到的最优解
			if slf.OptimizeUseless >= 3 && slf.OptimizeCount < slf.MaxOptimizeCount-10 {
				slf.OptimizeUseless = 0

				// 刷新禁排课位
				Refresh_FORBIDDEN_STAGE_CHANGE_WEEK(slf)
				Refresh_FORBIDDEN_STAGE_CHANGE_WEEK_NOSAME(slf)

				// 处理禁排位置课程
				slf.ClearForbid()

				addNum := len(slf.Waits) - waitNumber
				if addNum < 0 {
					addNum = 0
				}
				slf.OptimizeCount -= addNum
				if slf.OptimizeCount < 0 {
					slf.OptimizeCount = 0
				}
				fmt.Println("到达极限优化，策略调整，附加额外操作次数 >>", addNum)
			}
		}
		slf.Optimize()
	} else {
		fmt.Println("已达到设定的最大优化次数，结束优化，当前剩余冲突：", len(slf.Waits))
		slf.Save()
		slf.IsRunning = false
	}

}

// 自动排课
func (slf *Plan) AutoBuild() {
	slf.IsRunning = true
	// 移除所有已排课程放到待排课程中
	for _, week := range slf.Journeys {
		for _, section := range week {
			for _, course := range section {
				slf.Waits = append(slf.Waits, course)
			}
		}
	}
	weekNum := len(slf.Journeys)
	sectionNum := len(slf.Journeys[0])
	slf.Journeys = [][][]*Course{}
	for i := 0; i < weekNum; i++ {
		slf.Journeys = append(slf.Journeys, [][]*Course{})
		for a := 0; a < sectionNum; a++ {
			slf.Journeys[i] = append(slf.Journeys[i], []*Course{})
		}
	}

	// 刷新禁排课位
	Refresh_FORBIDDEN_STAGE_CHANGE_WEEK(slf)
	Refresh_FORBIDDEN_STAGE_CHANGE_WEEK_NOSAME(slf)

	// 优先级重排序
	var frist []*Course
	var other []*Course
	for _, wait := range slf.Waits {
		isFrist := false
		for _, s := range slf.MustContinuousTwo {
			if wait.Name == s {
				isFrist = true
				break
			}
		}
		if isFrist {
			frist = append(frist, wait)
		} else {
			other = append(other, wait)
		}
	}
	slf.Waits = append(frist, other...)

	// 排课结束后应该移除待排课的列表
	// 正常排课
	for _, course := range slf.Waits {
		// 存在可排课位随机抽取
		if allowSections := GetAllows(slf, course); len(allowSections) > 0 {
			target := allowSections[krand.Int(0, len(allowSections))]

			// 尽量保证每天特定课程不超过多少节
			if slf.WeekNumRangMax > 0 {
				var newAllows [][]int
				for _, ws := range allowSections {
					if GetWeeklyNumber(slf, course, ws[0]) <= slf.WeekNumRangMax-1 {
						newAllows = append(newAllows, ws)
					}
				}
				if len(newAllows) > 0 {
					allowSections = newAllows
				}
			}

			// 满足尽量连排，查找这个课程所在课位
			if slf.TryContinuous {
				placed := GetSection(slf, course)
				// 如果允许课程存在于已放置上或下课位  待定：且上上下下课位不存在则使用该课位
				for _, allow := range allowSections {
					for _, place := range placed {
						if allow[0] == place[0] {
							if allow[1] == place[1]-1 || allow[1] == place[1]+1 {
								target = allow
								break
							}
						}
					}
				}
			}

			course.IsRemoveWait = true
			slf.Journeys[target[0]][target[1]] = append(slf.Journeys[target[0]][target[1]], course)
		} else { // 无可排课位

		}
	}

	// 清理待排课区域
	var newWaits []*Course
	for _, course := range slf.Waits {
		if !course.IsRemoveWait {
			newWaits = append(newWaits, course)
		} else {
			course.IsRemoveWait = false
		}
	}
	slf.Waits = newWaits
	slf.Save()
	slf.IsRunning = false
}

// 保存排课方案
func (slf *Plan) Save() {
	// 如果已存在则删除
	filename := "assets/data/" + slf.Name + ".iscss"
	if _, err := os.Stat(filename); err == nil {
		if err = os.Remove(filename); err != nil {
			logs.Error(err)
		}
	}
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		logs.Error(err)
	} else {
		if data, err := json.Marshal(slf); err != nil {
			logs.Error(err)
		} else {
			_, err = file.Write(data)
			if err != nil {
				logs.Error(err)
			}
		}
	}
}

// 添加可用教室
func (slf *Plan) AddRoom(room string) {
	if strings.TrimSpace(room) != "" {
		for _, s := range slf.Rooms {
			if s == room {
				return
			}
		}
		slf.Rooms = append(slf.Rooms, room)
	}
	slf.Save()
}

// 获取当前课时数
func (slf *Plan) NowTotal() int {
	total := 0
	total = len(slf.Waits)
	for _, week := range slf.Journeys {
		for _, section := range week {
			total += len(section)
		}
	}
	slf.NowCourseTotal = total
	return total
}

func (slf *Plan) Draw() error {
	if err := slf.draw(""); err != nil {
		return err
	}
	for _, s := range slf.Stage {
		if err := slf.draw(s); err != nil {
			return err
		}
	}
	return nil
}

// 绘制最终课表到文件
func (slf *Plan) draw(stage string) error {
	if xlsx, err := xlsx.OpenFile("assets/temp/最终结果模板.xlsx"); err != nil {
		return err
	} else {
		sheet := xlsx.Sheets[0]
		// 渲染所需课位
		for w, week := range slf.Journeys {
			// 增加周次标题
			cell := sheet.Rows[0].Cells[w+1]
			cell.SetString(fmt.Sprint("天", w+1))

			for s, _ := range week {
				// 增加节次标题
				if w == 0 {
					row := sheet.AddRow()
					title := row.AddCell()
					title.SetString(fmt.Sprint("第", s+1, "节"))
					style := title.GetStyle()
					style.Font.Bold = true
					style.Font.Size = 10
					title.SetStyle(style)
					for i := 0; i < len(slf.Journeys); i++ {
						row.AddCell().SetString("")
					}
				}
			}
		}

		// 课位渲染
		for w, week := range slf.Journeys {
			for s, section := range week {
				for _, course := range section {
					if course.Stage == stage || stage == "" {
						cell := sheet.Rows[s+1].Cells[w+1]
						cell.SetString(fmt.Sprint(
							cell.String(),
							"\r\n",
							course.Name, "|", course.TeacherStr, "|", course.Room,
						))

						style := cell.GetStyle()
						style.Font.Size = 10
						cell.SetStyle(style)
					}
				}
			}
		}

		for _, row := range sheet.Rows {
			row.Hidden = false
			for _, cell := range row.Cells {
				style := cell.GetStyle()
				style.Alignment.WrapText = true
				cell.SetStyle(style)
			}
		}

		for _, row := range xlsx.Sheets[1].Rows {
			for _, cell := range row.Cells {
				sty := cell.GetStyle()
				sty.Font.Size = 10
				cell.SetStyle(sty)
			}
		}

		if err := os.MkdirAll(slf.Name, os.ModeDir); err != nil {
			return errors.New("the zips directory could not be initialized\r\n" + err.Error())
		}

		var err error
		if stage == "" {
			err = xlsx.Save(slf.Name + "/" + slf.Name + "_总课表.xlsx")
		} else {
			err = xlsx.Save(slf.Name + "/" + slf.Name + "_" + stage + ".xlsx")
		}
		if err != nil {
			fmt.Printf(err.Error())
		}
	}
	return nil
}

const (
	FORBIDDEN_DEFAULT                  = "[周, 节, 类型] 默认全年级统一禁排特定课位"
	FORBIDDEN_STAGE                    = "[周, 节, 类型, 年级] 这个课位特定年级将禁排"
	FORBIDDEN_STAGE_CHANGE_WEEK        = "[周, 当前随机节, 类型, 年级] 这个年级将随机禁排某周任一课位"
	FORBIDDEN_STAGE_CHANGE_WEEK_NOSAME = "[周, 当前随机节, 类型, 年级] 这个年级将随机禁排某周任一课位，且不与其他年级重复"
	FORBIDDEN_TEACHER                  = "[周, 节, 类型, 老师] 这个老师特定课位哪节不允许排课"
)

// 增加禁排课位
func (slf *Plan) AddForbidden(week, section int, other ...string) {
	week = week - 1
	section = section - 1

	// 无任何其他参数则全年级统一默认禁排
	if len(other) == 0 {
		slf.Forbidden = append(slf.Forbidden, []interface{}{week, section, FORBIDDEN_DEFAULT})
	} else {
		switch other[0] {
		case FORBIDDEN_DEFAULT:
			slf.Forbidden = append(slf.Forbidden, []interface{}{week, section, FORBIDDEN_DEFAULT})
		case FORBIDDEN_STAGE:
			slf.Forbidden = append(slf.Forbidden, []interface{}{week, section, FORBIDDEN_STAGE, other[1]})
		case FORBIDDEN_STAGE_CHANGE_WEEK:
			slf.Forbidden = append(slf.Forbidden, []interface{}{week, section, FORBIDDEN_STAGE_CHANGE_WEEK, other[1]})
		case FORBIDDEN_STAGE_CHANGE_WEEK_NOSAME:
			slf.Forbidden = append(slf.Forbidden, []interface{}{week, section, FORBIDDEN_STAGE_CHANGE_WEEK_NOSAME, other[1]})
		case FORBIDDEN_TEACHER:
			slf.Forbidden = append(slf.Forbidden, []interface{}{week, section, FORBIDDEN_TEACHER, other[1]})
		}
	}
	slf.Save()
}

// 禁排课位课程至待排区
func (slf *Plan) ClearForbid() {
	for _, forbidden := range slf.Forbidden {
		week := getForbiddenNumber(forbidden[0])
		section := getForbiddenNumber(forbidden[1])
		switch forbidden[2].(string) {
		case FORBIDDEN_DEFAULT:
			slf.Waits = append(slf.Waits, slf.Journeys[week][section]...)
			slf.Journeys[week][section] = []*Course{}
		case FORBIDDEN_STAGE, FORBIDDEN_STAGE_CHANGE_WEEK, FORBIDDEN_STAGE_CHANGE_WEEK_NOSAME:
			stage := forbidden[3].(string)
			var newSections []*Course
			for _, course := range slf.Journeys[week][section] {
				if course.Stage == stage {
					slf.Waits = append(slf.Waits, course)
				} else {
					newSections = append(newSections, course)
				}
			}
			slf.Journeys[week][section] = newSections
		case FORBIDDEN_TEACHER:
			teacher := forbidden[3].(string)
			var newSections []*Course
			for _, course := range slf.Journeys[week][section] {
				for _, t := range course.Teachers {
					if t == teacher {
						slf.Waits = append(slf.Waits, course)
						break
					} else {
						newSections = append(newSections, course)
					}
				}
			}
			slf.Journeys[week][section] = newSections

		}
	}
	slf.Save()
}
