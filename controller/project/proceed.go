package project

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

// Proceed Project
// @postfilter("Boss")
func (project *Project) Proceed(w http.ResponseWriter, r *http.Request) {
	fun := "project.Project.Proceed"

	var startReq *ProceedReq
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

	if projectAPI.State != api.StatePause {
		detail := "state must be paused when proceed"
		clog.Error("%s param err: %v, req: %v", fun, detail, startReq)
		project.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	if _, err := service.NewProject().Proceed(startReq.ID); err != nil {
		clog.Error("%s project.Proceed err: %v, req: %v", fun, err, startReq)
		project.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &ProceedRsp{}
	project.ReplyOk(w, resp)

	return
}
