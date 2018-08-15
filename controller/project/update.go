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

type UpdateReq struct {
	ID          int64           `json:"id"`
	CID         int64           `json:"cid"`
	Disabled    bool            `json:"disabled"`
	TargetTotal int             `json:"target_total"`
	WorkerTotal int             `json:"worker_total"`
	Msg         json.RawMessage `json:"msg,omitempty"`
	QPS         int             `json:"qps"`
}

func (updateReq *UpdateReq) Regular() (ok bool) {
	fun := "project.UpdateReq.Regular"

	if updateReq == nil {
		return
	}

	if updateReq.ID <= 0 {
		return
	}

	if updateReq.TargetTotal <= 0 {
		return
	}

	if updateReq.WorkerTotal <= 0 {
		return
	}

	if updateReq.QPS < 0 {
		return
	}

	cid := updateReq.CID
	constraintAPI, err := service.NewConstraint().Get(cid)
	if err != nil {
		clog.Error("%s constraint.Get error: %v, req: %v", fun, err, cid)
		return
	}
	if constraintAPI == nil || constraintAPI.Disabled {
		return
	}

	if len(updateReq.Msg) == 0 {
		return
	}

	ok = true
	return
}

type UpdateRsp struct {
	Project *api.Project `json:"project"`
}

// Update Project
// @postfilter("Boss")
func (project *Project) Update(w http.ResponseWriter, r *http.Request) {
	fun := "project.Project.Update"

	var updateReq *UpdateReq
	if err := json.Unmarshal(project.ReadBody(r), &updateReq); err != nil || !updateReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, updateReq)
		project.ReplyFail(w, lib.CodePara)
		return
	}

	projectAPI, err := service.NewProject().Get(updateReq.ID)
	if err != nil {
		clog.Error("%s project.Get error: %v, req: %v", fun, err, updateReq)
		project.ReplyFail(w, lib.CodeSrv)
		return
	}

	if projectAPI == nil {
		detail := api.CodeMap[api.CodeProjectNotExist]
		clog.Error("%s param err: %v, req: %v", fun, detail, updateReq)
		project.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	if projectAPI.State == api.StateStarted {
		detail := "no modify when state is started"
		clog.Error("%s param err: %v, req: %v", fun, detail, updateReq)
		project.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	projectAPI.Disabled = updateReq.Disabled
	projectAPI.TargetTotal = updateReq.TargetTotal
	projectAPI.WorkerTotal = updateReq.WorkerTotal
	projectAPI.CID = updateReq.CID
	projectAPI.Msg = updateReq.Msg
	projectAPI.Ut = time.Now().Unix()
	projectAPI.QPS = updateReq.QPS

	if err := service.NewProject().Upsert(projectAPI); err != nil {
		clog.Error("%s project.Upsert err: %v, req: %v", fun, err, updateReq)
		project.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &UpdateRsp{
		Project: projectAPI,
	}
	project.ReplyOk(w, resp)

	return
}
