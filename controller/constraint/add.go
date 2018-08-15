package constraint

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/simplejia/clog"
	"lib"
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/service"
	
)

type AddReq struct {
	TargetTotal   int                    `json:"target_total"`
	WorkerTotal   int                    `json:"worker_total"`
	ConditionKind api.ConditionKind `json:"condition_kind"`
	Condition     json.RawMessage        `json:"condition"`
	QPS           int                    `json:"qps"`
}

func (addReq *AddReq) Regular() (ok bool) {
	if addReq == nil {
		return
	}

	if addReq.TargetTotal <= 0 {
		return
	}

	if addReq.WorkerTotal <= 0 {
		return
	}

	if addReq.QPS < 0 {
		return
	}

	constraintAPI := &api.Constraint{
		ConditionKind: addReq.ConditionKind,
		Condition:     addReq.Condition,
		TargetTotal:   addReq.TargetTotal,
		WorkerTotal:   addReq.WorkerTotal,
	}
	if !constraintAPI.RegularCondition() {
		return
	}

	ok = true
	return
}

type AddRsp struct {
	Constraint *api.Constraint `json:"constraint"`
}

// Add Constraint
// @postfilter("Boss")
func (constraint *Constraint) Add(w http.ResponseWriter, r *http.Request) {
	fun := "constraint.Constraint.Add"

	var addReq *AddReq
	if err := json.Unmarshal(constraint.ReadBody(r), &addReq); err != nil || !addReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, addReq)
		constraint.ReplyFail(w, lib.CodePara)
		return
	}

	constraintAPI := api.NewConstraint()
	constraintAPI.State = api.StateReady
	constraintAPI.TargetTotal = addReq.TargetTotal
	constraintAPI.WorkerTotal = addReq.WorkerTotal
	constraintAPI.ConditionKind = addReq.ConditionKind
	constraintAPI.Condition = addReq.Condition
	constraintAPI.QPS = addReq.QPS
	constraintAPI.Ct = time.Now().Unix()

	if err := service.NewConstraint().Upsert(constraintAPI); err != nil {
		clog.Error("%s constraint.Upsert err: %v, req: %v", fun, err, constraintAPI)
		constraint.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &AddRsp{
		Constraint: constraintAPI,
	}
	constraint.ReplyOk(w, resp)

	return
}
