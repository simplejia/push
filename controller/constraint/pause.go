package constraint

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/clog/api"
	"lib"
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/service"
	
)

type PauseReq struct {
	ID int64 `json:"id"`
}

func (pauseReq *PauseReq) Regular() (ok bool) {
	if pauseReq == nil {
		return
	}

	if pauseReq.ID <= 0 {
		return
	}

	ok = true
	return
}

type PauseRsp struct {
}

// Pause Constraint
// @postfilter("Boss")
func (constraint *Constraint) Pause(w http.ResponseWriter, r *http.Request) {
	fun := "constraint.Constraint.Pause"

	var pauseReq *PauseReq
	if err := json.Unmarshal(constraint.ReadBody(r), &pauseReq); err != nil || !pauseReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, pauseReq)
		constraint.ReplyFail(w, lib.CodePara)
		return
	}

	constraintAPI, err := service.NewConstraint().Get(pauseReq.ID)
	if err != nil {
		clog.Error("%s constraint.Get error: %v, req: %v", fun, err, pauseReq)
		constraint.ReplyFail(w, lib.CodeSrv)
		return
	}

	if constraintAPI == nil || constraintAPI.Disabled {
		detail := api.CodeMap[api.CodeConstraintNotExist]
		clog.Error("%s param err: %v, req: %v", fun, detail, pauseReq)
		constraint.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	if constraintAPI.State != api.StateStarted {
		detail := "state must be started when paused"
		clog.Error("%s param err: %v, req: %v", fun, detail, pauseReq)
		constraint.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	if _, err := service.NewConstraint().Pause(pauseReq.ID); err != nil {
		clog.Error("%s constraint.Pause err: %v, req: %v", fun, err, pauseReq)
		constraint.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &PauseRsp{}
	constraint.ReplyOk(w, resp)

	return
}
