package constraint

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/clog"
	"lib"
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/service"
	
)

type ProceedReq struct {
	ID int64 `json:"id"`
}

func (startReq *ProceedReq) Regular() (ok bool) {
	if startReq == nil {
		return
	}

	if startReq.ID <= 0 {
		return
	}

	ok = true
	return
}

type ProceedRsp struct {
}

// Proceed Constraint
// @postfilter("Boss")
func (constraint *Constraint) Proceed(w http.ResponseWriter, r *http.Request) {
	fun := "constraint.Constraint.Proceed"

	var startReq *ProceedReq
	if err := json.Unmarshal(constraint.ReadBody(r), &startReq); err != nil || !startReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, startReq)
		constraint.ReplyFail(w, lib.CodePara)
		return
	}

	constraintAPI, err := service.NewConstraint().Get(startReq.ID)
	if err != nil {
		clog.Error("%s constraint.Get error: %v, req: %v", fun, err, startReq)
		constraint.ReplyFail(w, lib.CodeSrv)
		return
	}

	if constraintAPI == nil || constraintAPI.Disabled {
		detail := api.CodeMap[api.CodeConstraintNotExist]
		clog.Error("%s param err: %v, req: %v", fun, detail, startReq)
		constraint.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	if constraintAPI.State != api.StatePause {
		detail := "state must be paused when proceed"
		clog.Error("%s param err: %v, req: %v", fun, detail, startReq)
		constraint.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	if _, err := service.NewConstraint().Proceed(startReq.ID); err != nil {
		clog.Error("%s constraint.Proceed err: %v, req: %v", fun, err, startReq)
		constraint.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &ProceedRsp{}
	constraint.ReplyOk(w, resp)

	return
}
