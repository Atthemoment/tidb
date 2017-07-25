// Copyright 2017 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package plan

import (
	"fmt"
)

func setParents4FinalPlan(plan PhysicalPlan) {
	allPlans := []PhysicalPlan{plan}
	planMark := map[string]bool{}
	planMark[plan.ID()] = true
	for pID := 0; pID < len(allPlans); pID++ {
		allPlans[pID].SetParents()
		switch copPlan := allPlans[pID].(type) {
		case *PhysicalTableReader:
			setParents4FinalPlan(copPlan.tablePlan)
		case *PhysicalIndexReader:
			setParents4FinalPlan(copPlan.indexPlan)
		case *PhysicalIndexLookUpReader:
			setParents4FinalPlan(copPlan.indexPlan)
			setParents4FinalPlan(copPlan.tablePlan)
		}
		for _, p := range allPlans[pID].Children() {
			if !planMark[p.ID()] {
				allPlans = append(allPlans, p.(PhysicalPlan))
				planMark[p.ID()] = true
			}
		}
	}

	allPlans = allPlans[0:1]
	planMark[plan.ID()] = false
	for pID := 0; pID < len(allPlans); pID++ {
		for _, p := range allPlans[pID].Children() {
			p.AddParent(allPlans[pID])
			if planMark[p.ID()] {
				planMark[p.ID()] = false
				allPlans = append(allPlans, p.(PhysicalPlan))
			}
		}
	}
}

// ExplainInfo implements PhysicalPlan interface.
func (p *Limit) ExplainInfo() string {
	return fmt.Sprintf("offset:%v, count:%v", p.Offset, p.Count)
}
