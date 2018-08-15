package constraint

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"lib"

	"github.com/simplejia/clog"
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/service"
)

type LoadFromFileReq struct {
	Raw string `json:"raw"` // base64 encoding
}

func (loadFromFileReq *LoadFromFileReq) Regular() (ok bool) {
	if loadFromFileReq == nil {
		return
	}

	if loadFromFileReq.Raw == "" {
		return
	}

	ok = true
	return
}

type LoadFromFileRsp struct {
	Constraint *api.Constraint `json:"constraint"`
}

// LoadFromFile Constraint
// @postfilter("Boss")
func (constraint *Constraint) LoadFromFile(w http.ResponseWriter, r *http.Request) {
	fun := "constraint.Constraint.LoadFromFile"

	var loadFromFileReq *LoadFromFileReq
	if err := json.Unmarshal(constraint.ReadBody(r), &loadFromFileReq); err != nil || !loadFromFileReq.Regular() {
		clog.Error("%s param err: %v, req: %v", fun, err, loadFromFileReq)
		constraint.ReplyFail(w, lib.CodePara)
		return
	}

	data, err := base64.StdEncoding.DecodeString(loadFromFileReq.Raw)
	if err != nil {
		detail := "raw encoding unexpected"
		clog.Error("%s DecodeString err: %v, req: %v", fun, err, loadFromFileReq.Raw)
		constraint.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	mids := []int64{}

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		mid, err := strconv.ParseInt(line, 10, 64)
		if err != nil || mid <= 0 {
			clog.Error("%s parse err: %v, req: %v", fun, err, line)
			continue
		}

		mids = append(mids, mid)
	}

	if err := scanner.Err(); err != nil {
		clog.Error("%s scan err: %v, req: %v", fun, err, loadFromFileReq)
		constraint.ReplyFail(w, lib.CodePara)
		return
	}

	if len(mids) == 0 {
		detail := "no mids added"
		clog.Error("%s err: %v, req: %v", fun, detail, loadFromFileReq)
		constraint.ReplyFailWithDetail(w, lib.CodePara, detail)
		return
	}

	constraintAPI := api.NewConstraint()
	constraintAPI.State = api.StateReady
	constraintAPI.TargetTotal = len(mids)
	constraintAPI.WorkerTotal = 1
	constraintAPI.ConditionKind = api.ConditionKindConcrete
	conditionConcrete := &api.ConditionConcrete{
		Mids: mids,
	}
	bs, _ := json.Marshal(conditionConcrete)
	constraintAPI.Condition = bs
	constraintAPI.Ct = time.Now().Unix()

	if err := service.NewConstraint().Upsert(constraintAPI); err != nil {
		clog.Error("%s constraint.Upsert err: %v, req: %v", fun, err, constraintAPI)
		constraint.ReplyFail(w, lib.CodeSrv)
		return
	}

	resp := &LoadFromFileRsp{
		Constraint: constraintAPI,
	}
	constraint.ReplyOk(w, resp)

	return
}
