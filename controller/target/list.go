package target

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/clog/api"
	"lib"
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/service"
)

type ListReq struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (listReq *ListReq) Regular() (ok bool) {
	if listReq == nil {
		return
	}

	if listReq.Limit <= 0 {
		listReq.Limit = 20
	}

	ok = true
	return
}

type ListRsp struct {
	List   []*api.Target `json:"list"`
	Offset int           `json:"offset"`
	Total  int           `json:"total"`
}

// @postfilter("Boss")
func (target *Target) List(w http.ResponseWriter, r *http.Request) {
	fun := "target.Target.List"

	var listReq *ListReq
	if err := json.Unmarshal(target.ReadBody(r), &listReq); err != nil || !listReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, listReq)
		target.ReplyFail(w, lib.CodePara)
		return
	}

	total, err := service.NewTarget().Count()
	if err != nil {
		clog.Error("%s target.Count err: %v, req: %v", fun, err, listReq)
		target.ReplyFail(w, lib.CodeSrv)
		return
	}

	targets, err := service.NewTarget().List(listReq.Offset, listReq.Limit)
	if err != nil {
		clog.Error("%s target.List err: %v, req: %v", fun, err, listReq)
		target.ReplyFail(w, lib.CodeSrv)
		return
	}

	list := targets

	resp := &ListRsp{
		List:   list,
		Offset: listReq.Offset + len(list),
		Total:  total,
	}
	target.ReplyOk(w, resp)

	return

}
