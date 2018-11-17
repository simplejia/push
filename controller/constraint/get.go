package constraint

import (
	"encoding/json"
	"net/http"

	"github.com/simplejia/clog/api"
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
	Constraint *api.Constraint `json:"constraint,omitempty"`
}

// Get 通过ID提取一个Constraint
// @postfilter("Boss")
func (constraint *Constraint) Get(w http.ResponseWriter, r *http.Request) {
	fun := "constraint.Constraint.Get"

	var getReq *GetReq
	if err := json.Unmarshal(constraint.ReadBody(r), &getReq); err != nil || !getReq.Regular() {
		clog.Error("%s param error: %v, req: %v", fun, err, getReq)
		constraint.ReplyFail(w, lib.CodePara)
		return
	}

	constraintAPI, err := service.NewConstraint().Get(getReq.ID)
	if err != nil {
		clog.Error("%s constraint.Get error: %v, req: %v", fun, err, getReq)
		constraint.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &GetRsp{
		Constraint: constraintAPI,
	}
	constraint.ReplyOk(w, resp)

	return
}
