package project

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/clog"
	"lib"
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/service"
)

type GetReq struct {
	ID int64 `json:"id"`
}

func (getReq *GetReq) Regular() (ok bool) {
	if getReq == nil {
		return
	}

	if getReq.ID <= 0 {
		return
	}

	ok = true
	return
}

type GetRsp struct {
	Project *api.Project `json:"project,omitempty"`
}

// @postfilter("Boss")
func (project *Project) Get(w http.ResponseWriter, r *http.Request) {
	fun := "project.Project.Get"

	var getReq *GetReq
	if err := json.Unmarshal(project.ReadBody(r), &getReq); err != nil || !getReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, getReq)
		project.ReplyFail(w, lib.CodePara)
		return
	}

	projectAPI, err := service.NewProject().Get(getReq.ID)
	if err != nil {
		clog.Error("%s project.Get err: %v, req: %v", fun, err, getReq)
		project.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &GetRsp{
		Project: projectAPI,
	}
	project.ReplyOk(w, resp)

	return
}
