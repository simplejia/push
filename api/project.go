package api

import (
	"encoding/json"
)

type Project struct {
	ID            int64            `json:"id" bson:"_id"`
	CID           int64            `json:"cid" bson:"cid"`
	State         State            `json:"state" bson:"state"` // no empty
	QPS           int              `json:"qps,omitempty" bson:"qps,omitempty"`
	Ct            int64            `json:"ct" bson"ct"`
	Ut            int64            `json:"ut" bson"ut"`
	Timestamps    map[string]int64 `json:"timestamps" bson:"timestamps"`
	Disabled      bool             `json:"disabled" bson:"disabled"`
	WorkerCounter int              `json:"worker_counter" bson:"worker_counter"` // no empty
	WorkerTotal   int              `json:"worker_total" bson:"worker_total"`
	TargetCounter int              `json:"target_counter" bson:"target_counter"` // no empty
	TargetTotal   int              `json:"target_total" bson:"target_total"`
	Msg           json.RawMessage  `json:"msg,omitempty" bson:"msg,omitempty"`
}

func NewProject() *Project {
	return &Project{}
}
