package project

import (
	"encoding/json"
	"net/http"

	"lib"

	"github.com/simplejia/clog/api"
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
	List   []*api.Project `json:"list"`
	Offset int            `json:"offset"`
	Total  int            `json:"total"`
}

// @postfilter("Boss")
func (project *Project) List(w http.ResponseWriter, r *http.Request) {
	fun := "project.Project.List"

	var listReq *ListReq
	if err := json.Unmarshal(project.ReadBody(r), &listReq); err != nil || !listReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, listReq)
		project.ReplyFail(w, lib.CodePara)
		return
	}

	total, err := service.NewProject().Count()
	if err != nil {
		clog.Error("%s project.Count err: %v, req: %v", fun, err, listReq)
		project.ReplyFail(w, lib.CodeSrv)
		return
	}

	projects, err := service.NewProject().List(listReq.Offset, listReq.Limit)
	if err != nil {
		clog.Error("%s project.List err: %v, req: %v", fun, err, listReq)
		project.ReplyFail(w, lib.CodeSrv)
		return
	}

	list := projects

	resp := &ListRsp{
		List:   list,
		Offset: listReq.Offset + len(list),
		Total:  total,
	}
	project.ReplyOk(w, resp)

	return

}
