package constraint

import (
	"encoding/json"
	"lib"
	"net/http"

	"github.com/simplejia/clog"
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
	List   []*api.Constraint `json:"list"`
	Offset int               `json:"offset"`
	Total  int               `json:"total"`
}

// @postfilter("Boss")
func (constraint *Constraint) List(w http.ResponseWriter, r *http.Request) {
	fun := "constraint.Constraint.List"

	var listReq *ListReq
	if err := json.Unmarshal(constraint.ReadBody(r), &listReq); err != nil || !listReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, listReq)
		constraint.ReplyFail(w, lib.CodePara)
		return
	}

	total, err := service.NewConstraint().Count()
	if err != nil {
		clog.Error("%s constraint.Count err: %v, req: %v", fun, err, listReq)
		constraint.ReplyFail(w, lib.CodeSrv)
		return
	}

	constraints, err := service.NewConstraint().List(listReq.Offset, listReq.Limit)
	if err != nil {
		clog.Error("%s constraint.List err: %v, req: %v", fun, err, listReq)
		constraint.ReplyFail(w, lib.CodeSrv)
		return
	}

	list := constraints

	resp := &ListRsp{
		List:   list,
		Offset: listReq.Offset + len(list),
		Total:  total,
	}
	constraint.ReplyOk(w, resp)

	return

}
