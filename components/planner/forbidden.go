package planner

import (
	"github.com/kercylan98/dev-kits/utils/krand"
)

// 刷新禁排周
func Refresh_FORBIDDEN_STAGE_CHANGE_WEEK(plan *Plan) {
	var newForbidden [][]interface{}
	for _, forbidden := range plan.Forbidden {
		if forbidden[2].(string) == FORBIDDEN_STAGE_CHANGE_WEEK {
			week := getForbiddenNumber(forbidden[0])
		retry:
			{
			}
			section := krand.Int(0, len(plan.Journeys[0]))
			// 检查规则是否重复
			for _, checkF := range plan.Forbidden {
				if checkF[2].(string) == FORBIDDEN_STAGE_CHANGE_WEEK && checkF[3].(string) == forbidden[3].(string) {
					if week == getForbiddenNumber(forbidden[0]) && section == getForbiddenNumber(forbidden[1]) {
						goto retry
					}
				}
			}
			f := []interface{}{
				forbidden[0],
				section,
				forbidden[2],
				forbidden[3],
			}
			newForbidden = append(newForbidden, f)
		} else {
			newForbidden = append(newForbidden, forbidden)
		}
	}
	plan.Forbidden = newForbidden
}

// 刷新禁排周
func Refresh_FORBIDDEN_STAGE_CHANGE_WEEK_NOSAME(plan *Plan) {
	var newForbidden [][]interface{}
	for _, forbidden := range plan.Forbidden {
		if forbidden[2].(string) == FORBIDDEN_STAGE_CHANGE_WEEK_NOSAME {
			week := getForbiddenNumber(forbidden[0])
		retry:
			{
			}
			section := krand.Int(0, len(plan.Journeys[0]))
			// 检查规则是否重复
			for _, checkF := range plan.Forbidden {
				if checkF[2].(string) == FORBIDDEN_STAGE_CHANGE_WEEK_NOSAME {
					if week == getForbiddenNumber(forbidden[0]) && section == getForbiddenNumber(forbidden[1]) {
						goto retry
					}
				}
			}
			f := []interface{}{
				forbidden[0],
				section,
				forbidden[2],
				forbidden[3],
			}
			newForbidden = append(newForbidden, f)
		} else {
			newForbidden = append(newForbidden, forbidden)
		}
	}
	plan.Forbidden = newForbidden
}
