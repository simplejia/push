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

type UpdateReq struct {
	ID            int64                  `json:"id"`
	Disabled      bool                   `json:"disabled"`
	TargetTotal   int                    `json:"target_total"`
	WorkerTotal   int                    `json:"worker_total"`
	ConditionKind api.ConditionKind `json:"condition_kind"`
	Condition     json.RawMessage        `json:"condition"`
	QPS           int                    `json:"qps"`
}

func (updateReq *UpdateReq) Regular() (ok bool) {
	if updateReq == nil {
		return
	}

	if updateReq.ID <= 0 {
		return
	}

	if updateReq.TargetTotal <= 0 {
		return
	}

	if updateReq.WorkerTotal <= 0 {
		return
	}

	if updateReq.QPS < 0 {
		return
	}

	constraintAPI := &api.Constraint{
		ConditionKind: updateReq.ConditionKind,
		Condition:     updateReq.Condition,
		TargetTotal:   updateReq.TargetTotal,
		WorkerTotal:   updateReq.WorkerTotal,
	}
	if !constraintAPI.RegularCondition() {
		return
	}

	ok = true
	return
}

type UpdateRsp struct {
	Constraint *api.Constraint `json:"constraint"`
}

// Update Constraint
// @postfilter("Boss")
func (constraint *Constraint) Update(w http.ResponseWriter, r *http.Request) {
	fun := "constraint.Constraint.Update"

	var updateReq *UpdateReq
	if err := json.Unmarshal(constraint.ReadBody(r), &updateReq); err != nil || !updateReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, updateReq)
		constraint.ReplyFail(w, lib.CodePara)
		return
	}

	constraintAPI, err := service.NewConstraint().Get(updateReq.ID)
	if err != nil {
		clog.Error("%s constraint.Get error: %v, req: %v", fun, err, updateReq)
		constraint.ReplyFail(w, lib.CodeSrv)
		return
	}

	if constraintAPI == nil {
		detail := api.CodeMap[api.CodeConstraintNotExist]
		clog.Error("%s param err: %v, req: %v", fun, detail, updateReq)
		constraint.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	if constraintAPI.State == api.StateStarted {
		detail := "no modify when state is started"
		clog.Error("%s param err: %v, req: %v", fun, detail, updateReq)
		constraint.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	constraintAPI.Disabled = updateReq.Disabled
	constraintAPI.TargetTotal = updateReq.TargetTotal
	constraintAPI.WorkerTotal = updateReq.WorkerTotal
	constraintAPI.ConditionKind = updateReq.ConditionKind
	constraintAPI.Condition = updateReq.Condition
	constraintAPI.QPS = updateReq.QPS
	constraintAPI.Ut = time.Now().Unix()

	if err := service.NewConstraint().Upsert(constraintAPI); err != nil {
		clog.Error("%s constraint.Upsert err: %v, req: %v", fun, err, updateReq)
		constraint.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &UpdateRsp{
		Constraint: constraintAPI,
	}
	constraint.ReplyOk(w, resp)

	return
}
