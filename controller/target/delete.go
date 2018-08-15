package target

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/clog"
	"lib"
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/service"
)

type DeleteReq struct {
	ID string `json:"id"`
}

func (deleteReq *DeleteReq) Regular() (ok bool) {
	if deleteReq == nil {
		return
	}

	if deleteReq.ID == "" {
		return
	}

	ok = true
	return
}

type DeleteRsp struct {
	Target *api.Target `json:"target"`
}

// Delete Target
// @postfilter("Boss")
func (target *Target) Delete(w http.ResponseWriter, r *http.Request) {
	fun := "target.Target.Delete"

	var deleteReq *DeleteReq
	if err := json.Unmarshal(target.ReadBody(r), &deleteReq); err != nil || !deleteReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, deleteReq)
		target.ReplyFail(w, lib.CodePara)
		return
	}

	if err := service.NewTarget().Delete(deleteReq.ID); err != nil {
		clog.Error("%s target.Delete err: %v, req: %v", fun, err, deleteReq)
		target.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &DeleteRsp{}
	target.ReplyOk(w, resp)

	return
}
