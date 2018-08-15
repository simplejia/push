package project

import (
	"encoding/json"
	"net/http"
	"time"

	"lib"

	"github.com/simplejia/clog"
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/service"
)

type AddReq struct {
	CID         int64           `json:"cid"`
	Msg         json.RawMessage `json:"msg,omitempty"`
	TargetTotal int             `json:"target_total"`
	WorkerTotal int             `json:"worker_total"`
	QPS         int             `json:"qps"`
}

func (addReq *AddReq) Regular() (ok bool) {
	fun := "project.AddReq.Regular"

	if addReq == nil {
		return
	}

	cid := addReq.CID
	constraintAPI, err := service.NewConstraint().Get(cid)
	if err != nil {
		clog.Error("%s constraint.Get error: %v, req: %v", fun, err, cid)
		return
	}
	if constraintAPI == nil || constraintAPI.Disabled {
		return
	}

	if len(addReq.Msg) == 0 {
		return
	}

	if addReq.TargetTotal <= 0 {
		return
	}

	if addReq.WorkerTotal <= 0 {
		return
	}

	if addReq.QPS < 0 {
		return
	}

	ok = true
	return
}

type AddRsp struct {
	Project *api.Project `json:"project"`
}

// @postfilter("Boss")
func (project *Project) Add(w http.ResponseWriter, r *http.Request) {
	fun := "project.Project.Add"

	var addReq *AddReq
	if err := json.Unmarshal(project.ReadBody(r), &addReq); err != nil || !addReq.Regular() {
		clog.Error("%s param error: %v, req: %v", fun, err, addReq)
		project.ReplyFail(w, lib.CodePara)
		return
	}

	projectAPI := api.NewProject()
	projectAPI.State = api.StateReady
	projectAPI.CID = addReq.CID
	projectAPI.Msg = addReq.Msg
	projectAPI.TargetTotal = addReq.TargetTotal
	projectAPI.WorkerTotal = addReq.WorkerTotal
	projectAPI.Ct = time.Now().Unix()
	projectAPI.QPS = addReq.QPS

	if err := service.NewProject().Upsert(projectAPI); err != nil {
		clog.Error("%s project.Upsert err: %v, req %v", fun, err, addReq)
		project.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &AddRsp{
		Project: projectAPI,
	}
	project.ReplyOk(w, resp)

	return
}
