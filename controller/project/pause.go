package project

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/clog"
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

// Pause Project
// @postfilter("Boss")
func (project *Project) Pause(w http.ResponseWriter, r *http.Request) {
	fun := "project.Project.Pause"

	var pauseReq *PauseReq
	if err := json.Unmarshal(project.ReadBody(r), &pauseReq); err != nil || !pauseReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, pauseReq)
		project.ReplyFail(w, lib.CodePara)
		return
	}

	projectAPI, err := service.NewProject().Get(pauseReq.ID)
	if err != nil {
		clog.Error("%s project.Get error: %v, req: %v", fun, err, pauseReq)
		project.ReplyFail(w, lib.CodeSrv)
		return
	}

	if projectAPI == nil || projectAPI.Disabled {
		detail := api.CodeMap[api.CodeProjectNotExist]
		clog.Error("%s param err: %v, req: %v", fun, detail, pauseReq)
		project.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	if projectAPI.State != api.StateStarted {
		detail := "state must be started when paused"
		clog.Error("%s param err: %v, req: %v", fun, detail, pauseReq)
		project.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	if _, err := service.NewProject().Pause(pauseReq.ID); err != nil {
		clog.Error("%s project.Pause err: %v, req: %v", fun, err, pauseReq)
		project.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &PauseRsp{}
	project.ReplyOk(w, resp)

	return
}
