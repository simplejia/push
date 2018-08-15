package project

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

// Start Project
// @postfilter("Boss")
func (project *Project) Start(w http.ResponseWriter, r *http.Request) {
	fun := "project.Project.Start"

	var startReq *StartReq
	if err := json.Unmarshal(project.ReadBody(r), &startReq); err != nil || !startReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, startReq)
		project.ReplyFail(w, lib.CodePara)
		return
	}

	projectAPI, err := service.NewProject().Get(startReq.ID)
	if err != nil {
		clog.Error("%s project.Get error: %v, req: %v", fun, err, startReq)
		project.ReplyFail(w, lib.CodeSrv)
		return
	}

	if projectAPI == nil || projectAPI.Disabled {
		detail := api.CodeMap[api.CodeProjectNotExist]
		clog.Error("%s param err: %v, req: %v", fun, detail, startReq)
		project.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	if projectAPI.State == api.StateStarted {
		detail := "state is already started"
		clog.Error("%s param err: %v, req: %v", fun, detail, startReq)
		project.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	cid := projectAPI.CID
	constraintAPI, err := service.NewConstraint().Get(cid)
	if err != nil {
		clog.Error("%s constraint.Get error: %v, req: %v", fun, err, cid)
		project.ReplyFail(w, lib.CodeSrv)
		return
	}
	if constraintAPI == nil || constraintAPI.Disabled {
		detail := api.CodeMap[api.CodeConstraintNotExist]
		clog.Error("%s param err: %v, req: %v", fun, detail, cid)
		project.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	if _, err := service.NewProject().Start(startReq.ID); err != nil {
		clog.Error("%s project.Start err: %v, req: %v", fun, err, startReq)
		project.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &StartRsp{}
	project.ReplyOk(w, resp)

	return
}
