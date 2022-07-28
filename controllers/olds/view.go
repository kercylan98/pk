package olds

import (
	"github.com/kercylan98/dev-kits/utils/kstreams"
	"pk/common/basic"
	"pk/components/planner"
)

type ViewController struct {
	basic.Controller
}

// 返回主视图
func (slf *ViewController) View() {
	var plans []string
	kstreams.EachMapSort(planner.GetPlans(), func(name string, plan *planner.Plan) {
		plans = append(plans, name)
	})

	slf.Data["plan"] = planner.GetOnlinePlan(slf.Ctx.Request)
	slf.Data["plans"] = plans
	if plan := planner.GetOnlinePlan(slf.Ctx.Request); plan != nil {
		slf.Data["isRunning"] = plan.IsRunning
	} else {
		slf.Data["isRunning"] = false
	}
	slf.TplName = "index.html"
}
