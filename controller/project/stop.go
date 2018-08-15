package project

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/clog"
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

// Stop Project
// @postfilter("Boss")
func (project *Project) Stop(w http.ResponseWriter, r *http.Request) {
	fun := "project.Project.Stop"

	var stopReq *StopReq
	if err := json.Unmarshal(project.ReadBody(r), &stopReq); err != nil || !stopReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, stopReq)
		project.ReplyFail(w, lib.CodePara)
		return
	}

	projectAPI, err := service.NewProject().Get(stopReq.ID)
	if err != nil {
		clog.Error("%s project.Get error: %v, req: %v", fun, err, stopReq)
		project.ReplyFail(w, lib.CodeSrv)
		return
	}

	if projectAPI == nil || projectAPI.Disabled {
		detail := api.CodeMap[api.CodeProjectNotExist]
		clog.Error("%s param err: %v, req: %v", fun, detail, stopReq)
		project.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	if projectAPI.State != api.StateStarted {
		detail := "state must be started when stopped"
		clog.Error("%s param err: %v, req: %v", fun, detail, stopReq)
		project.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	if _, err := service.NewProject().Stop(stopReq.ID); err != nil {
		clog.Error("%s project.Stop err: %v, req: %v", fun, err, stopReq)
		project.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &StopRsp{}
	project.ReplyOk(w, resp)

	return
}
