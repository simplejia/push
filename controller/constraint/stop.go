package constraint

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/clog/api"
	"lib"
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/service"
	
)

type StopReq struct {
	ID int64 `json:"id"`
}

func (stopReq *StopReq) Regular() (ok bool) {
	if stopReq == nil {
		return
	}

	if stopReq.ID <= 0 {
		return
	}

	ok = true
	return
}

type StopRsp struct {
}

// Stop Constraint
// @postfilter("Boss")
func (constraint *Constraint) Stop(w http.ResponseWriter, r *http.Request) {
	fun := "constraint.Constraint.Stop"

	var stopReq *StopReq
	if err := json.Unmarshal(constraint.ReadBody(r), &stopReq); err != nil || !stopReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, stopReq)
		constraint.ReplyFail(w, lib.CodePara)
		return
	}

	constraintAPI, err := service.NewConstraint().Get(stopReq.ID)
	if err != nil {
		clog.Error("%s constraint.Get error: %v, req: %v", fun, err, stopReq)
		constraint.ReplyFail(w, lib.CodeSrv)
		return
	}

	if constraintAPI == nil || constraintAPI.Disabled {
		detail := api.CodeMap[api.CodeConstraintNotExist]
		clog.Error("%s param err: %v, req: %v", fun, detail, stopReq)
		constraint.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	if constraintAPI.State != api.StateStarted {
		detail := "state must be started when stopped"
		clog.Error("%s param err: %v, req: %v", fun, detail, stopReq)
		constraint.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	if _, err := service.NewConstraint().Stop(stopReq.ID); err != nil {
		clog.Error("%s constraint.Stop err: %v, req: %v", fun, err, stopReq)
		constraint.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &StopRsp{}
	constraint.ReplyOk(w, resp)

	return
}
