package planner

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/kercylan98/dev-kits/utils/kfile"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var this *planner

func init() {
	this = &planner{
		plans:      map[string]*Plan{},
		onlinePlan: map[string]string{},
	}

	this.loadLocalhostPlan()
}

// 方案管理模型
type planner struct {
	sync.RWMutex
	plans      map[string]*Plan  // 所有排课方案
	onlinePlan map[string]string // 当前在线的排课方案
}

// 设置正在打开的方案
func SetOnlinePlan(request *http.Request, plan *Plan) {
	logs.Debug("Set online plan => ", strings.Split(request.RemoteAddr, ":")[0], plan.Name)
	this.onlinePlan[strings.Split(request.RemoteAddr, ":")[0]] = plan.Name
}

// 获取正在打开的方案
func GetOnlinePlan(request *http.Request) *Plan {
	logs.Debug("Get online plan => ", strings.Split(request.RemoteAddr, ":")[0])
	plan := this.plans[this.onlinePlan[strings.Split(request.RemoteAddr, ":")[0]]]
	if plan != nil {
		plan.NowTotal()
	}
	return plan
}

// 加载本地方案
func (slf *planner) loadLocalhostPlan() {
	slf.Lock()
	slf.plans = map[string]*Plan{}
	files, err := ioutil.ReadDir("assets/data")
	if err != nil {
		panic(err)
	}

	for _, fileInfo := range files {
		if strings.HasSuffix(fileInfo.Name(), ".iscss") {
			data, err := kfile.ReadOnce("assets/data/" + fileInfo.Name())
			if err != nil {
				panic(err)
			}

			var plan = new(Plan)
			if err := json.Unmarshal(data, plan); err != nil {
				panic(err)
			}
			plan.IsRunning = false
			AddPlan(plan)
		}
	}
	slf.Unlock()
}

// 新增方案
func AddPlan(plan *Plan) {
	this.plans[plan.Name] = plan
	plan.Save()
}

// 获取所有方案
func GetPlans() map[string]*Plan {
	return this.plans
}

// 根据方案名称获取方案
func GetPlan(name string) *Plan {
	return this.plans[name]
}

// 根据模板地址加载方案数据
func LoadData(planName string, templatePath string) error {
	plan := GetPlan(planName)
	if plan == nil {
		return errors.New("plan not exist")
	}
	// 打开文件
	if xlsx, err := xlsx.OpenFile(templatePath); err != nil {
		return err
	} else {
		courseSheet := xlsx.Sheets[0]
		studentSheet := xlsx.Sheets[1]
		roomSheet := xlsx.Sheets[2]

		// 加载课程信息
		for ci, cr := range courseSheet.Rows {
			if ci > 0 {
				weeklyNumber, err := cr.Cells[4].Int()
				if err != nil {
					return err
				}

				for number := 0; number < weeklyNumber; number++ {
					plan.CourseCount++
					course := &Course{
						Number:       strings.TrimSpace(cr.Cells[0].String()),
						Name:         strings.TrimSpace(cr.Cells[1].String()),
						Teachers:     strings.Split(strings.TrimSpace(cr.Cells[3].String()), ","),
						WeeklyNumber: weeklyNumber,
						Students:     nil,
						Stage:        strings.TrimSpace(cr.Cells[6].String()),
						Room:         strings.TrimSpace(cr.Cells[5].String()),
						TeacherStr:   strings.TrimSpace(cr.Cells[3].String()),
					}

					// 存储阶段
					exist := false
					for _, s := range plan.Stage {
						if course.Stage == s {
							exist = true
							break
						}
					}
					if !exist {
						plan.Stage = append(plan.Stage, course.Stage)
					}

					// 加载课程学生
					var students []*Student
					for si, sr := range studentSheet.Rows {
						if si > 0 && strings.TrimSpace(sr.Cells[0].String()) == course.Name {
							student := &Student{
								Number: strings.TrimSpace(sr.Cells[2].String()),
								Name:   strings.TrimSpace(sr.Cells[3].String()),
							}
							students = append(students, student)
							// 添加学生对应课程缓存数据
							plan.StudentHasCourse[student.Number] = append(plan.StudentHasCourse[student.Number], course)
						}
					}

					course.Students = students

					// 同时将教室添加到全局可用教室
					plan.AddRoom(course.Room)

					// 添加课程到待排课程
					plan.Waits = append(plan.Waits, course)
				}

			}

		}

		// 加载全局教室
		for ri, rr := range roomSheet.Rows {
			if ri > 0 {
				plan.AddRoom(rr.Cells[0].String())
			}
		}
	}
	plan.Save()
	return nil
}