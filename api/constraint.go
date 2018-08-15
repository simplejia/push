package api

import (
	"encoding/json"
)

type Constraint struct {
	ID            int64            `json:"id" bson:"_id"`
	Disabled      bool             `json:"disabled" bson:"disabled"`
	QPS           int              `json:"qps,omitempty" bson:"qps,omitempty"`
	Ct            int64            `json:"ct" bson"ct"`
	Ut            int64            `json:"ut" bson"ut"`
	State         State            `json:"state" bson:"state"`
	Timestamps    map[string]int64 `json:"timestamps" bson:"timestamps"`
	TargetTotal   int              `json:"target_total" bson:"target_total"`
	TargetCounter int              `json:"target_counter" bson:"target_counter"`
	WorkerTotal   int              `json:"worker_total" bson:"worker_total"`
	WorkerCounter int              `json:"worker_counter" bson:"worker_counter"`
	ConditionKind ConditionKind    `json:"condition_kind" bson:"condition_kind"`
	Condition     json.RawMessage  `json:"condition,omitempty" bson:"condition,omitempty"`
}

func NewConstraint() *Constraint {
	return &Constraint{}
}

func (constraint *Constraint) RegularCondition() (ok bool) {
	switch constraint.ConditionKind {
	case ConditionKindConcrete:
		var conditionConcrete *ConditionConcrete
		if err := json.Unmarshal(constraint.Condition, &conditionConcrete); err != nil || !conditionConcrete.Regular() {
			return
		}
		if len(conditionConcrete.Mids) != constraint.TargetTotal {
			return
		}
	default:
		return
	}

	ok = true
	return
}

type ConditionConcrete struct {
	Mids []int64 `json:"mids,omitempty" bson:"mids,omitempty"`
}

func (conditionConcrete *ConditionConcrete) Regular() (ok bool) {
	if conditionConcrete == nil {
		return
	}

	if len(conditionConcrete.Mids) == 0 {
		return
	}

	ok = true
	return
}
