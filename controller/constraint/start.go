package constraint

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/clog"
	"lib"
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/service"
	
)

type StartReq struct {
	ID int64 `json:"id"`
}

func (startReq *StartReq) Regular() (ok bool) {
	if startReq == nil {
		return
	}

	if startReq.ID <= 0 {
		return
	}

	ok = true
	return
}

type StartRsp struct {
}

// Start Constraint
// @postfilter("Boss")
func (constraint *Constraint) Start(w http.ResponseWriter, r *http.Request) {
	fun := "constraint.Constraint.Start"

	var startReq *StartReq
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

	if constraintAPI.State == api.StateStarted {
		detail := "state is already started"
		clog.Error("%s param err: %v, req: %v", fun, detail, startReq)
		constraint.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	if _, err := service.NewConstraint().Start(startReq.ID); err != nil {
		clog.Error("%s constraint.Start err: %v, req: %v", fun, err, startReq)
		constraint.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &StartRsp{}
	constraint.ReplyOk(w, resp)

	return
}
