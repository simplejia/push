package target

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/clog/api"
	"lib"
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/service"
)

type ListByCidReq struct {
	Cid    int64 `json:"cid"`
	Offset int   `json:"offset"`
	Limit  int   `json:"limit"`
}

func (listByCidReq *ListByCidReq) Regular() (ok bool) {
	if listByCidReq == nil {
		return
	}

	if listByCidReq.Cid <= 0 {
		return
	}

	if listByCidReq.Limit <= 0 {
		listByCidReq.Limit = 20
	}

	ok = true
	return
}

type ListByCidRsp struct {
	List   []*api.Target `json:"list"`
	Offset int           `json:"offset"`
	Total  int           `json:"total"`
}

// @postfilter("Boss")
func (target *Target) ListByCid(w http.ResponseWriter, r *http.Request) {
	fun := "target.Target.ListByCid"

	var listByCidReq *ListByCidReq
	if err := json.Unmarshal(target.ReadBody(r), &listByCidReq); err != nil || !listByCidReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, listByCidReq)
		target.ReplyFail(w, lib.CodePara)
		return
	}

	total, err := service.NewTarget().CountByCid(listByCidReq.Cid)
	if err != nil {
		clog.Error("%s target.CountByCid err: %v, req: %v", fun, err, listByCidReq)
		target.ReplyFail(w, lib.CodeSrv)
		return
	}

	targets, err := service.NewTarget().ListByCid(listByCidReq.Cid, listByCidReq.Offset, listByCidReq.Limit)
	if err != nil {
		clog.Error("%s target.ListByCid err: %v, req: %v", fun, err, listByCidReq)
		target.ReplyFail(w, lib.CodeSrv)
		return
	}

	list := targets

	resp := &ListByCidRsp{
		List:   list,
		Offset: listByCidReq.Offset + len(list),
		Total:  total,
	}
	target.ReplyOk(w, resp)

	return

}
