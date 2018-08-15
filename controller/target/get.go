package target

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/clog"
	"lib"
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/service"
)

type GetReq struct {
	CID int64 `json:"cid"`
	MID int64 `json:"mid"`
}

func (getReq *GetReq) Regular() (ok bool) {
	if getReq == nil {
		return
	}

	if getReq.CID < 0 || getReq.MID <= 0 {
		return
	}

	ok = true
	return
}

type GetRsp struct {
	Target *api.Target `json:"target,omitempty"`
}

// @postfilter("Boss")
func (target *Target) Get(w http.ResponseWriter, r *http.Request) {
	fun := "target.Target.Get"

	var getReq *GetReq
	if err := json.Unmarshal(target.ReadBody(r), &getReq); err != nil || !getReq.Regular() {
		clog.Error("%s param err: %v, reg: %v", fun, err, getReq)
		target.ReplyFail(w, lib.CodePara)
		return
	}

	targetAPI, err := service.NewTarget().Get(getReq.CID, getReq.MID)
	if err != nil {
		clog.Error("%s target.Get err: %v, req: %v", fun, err, getReq)
		target.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &GetRsp{
		Target: targetAPI,
	}
	target.ReplyOk(w, resp)

	return
}
