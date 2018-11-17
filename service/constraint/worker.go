package constraint

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/simplejia/clog/api"
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/model"
	constraint_model "github.com/simplejia/push/model/constraint"
)

func CreateGetMembersFunc(
	constraintModel *constraint_model.Constraint,
	targetNum int,
) (f func() ([]int64, bool, error), err error) {
	conditionKind := constraintModel.ConditionKind

	switch conditionKind {
	case api.ConditionKindConcrete:
		conditionConcrete := &api.ConditionConcrete{}
		if err = json.Unmarshal(constraintModel.Condition, &conditionConcrete); err != nil {
			return
		}

		if conditionConcrete == nil {
			err = errors.New("condition invalid")
			return
		}

		firstRun := true
		f = func() (mids []int64, exhausted bool, err error) {
			if !firstRun {
				exhausted = true
				return
			}

			firstRun = false

			mids = conditionConcrete.Mids
			if len(mids) == 0 {
				return
			}

			return
		}
	default:
		err = fmt.Errorf("Unknown conditionKind: %v", conditionKind)
		return
	}
	return
}

func Worker(cid int64, wid string, channelFromMaster, channelToMaster chan *Resource) (err error) {
	fun := "constraint.Worker"
	clog.Info("%s worker started, req: %v,%v", fun, cid, wid)

	var signal api.State

	defer func() {
		channelToMaster <- &Resource{
			CID:    cid,
			WID:    wid,
			Signal: signal,
		}

		if err != nil {
			clog.Error("%s err: %v, req: %v,%v", fun, err, cid, wid)
		}
	}()

	heartbeatTicker := time.Tick(api.KeepaliveTimeThreshold)
	errCounter, errThreshold := 0, api.ErrorThreshold

	constraintModel := model.NewConstraint()
	constraintModel.ID = cid
	constraintModel, err = constraintModel.Get()
	if err != nil {
		signal = api.StateFailed
		return
	}

	if constraintModel == nil {
		err = errors.New("constraint.Get empty")
		signal = api.StateFailed
		return
	}

	qps, targetNum := constraintModel.QPS, 100
	if qps <= 0 {
		qps = 10000
	}

	workerTicker := time.Tick(time.Duration(int(time.Second) * targetNum * constraintModel.WorkerTotal / qps))

	getMembersFunc, err := CreateGetMembersFunc(constraintModel, targetNum)
	if err != nil {
		signal = api.StateFailed
		return
	}

	for {
		if errCounter >= errThreshold {
			err = errors.New("failed, err threshold exceeded")
			signal = api.StateFailed
			return
		}

		select {
		case <-heartbeatTicker: // send heartbeat timestamp.
			if err := constraintModel.SetTimestamp(wid, time.Now().Unix()); err != nil {
				errCounter++
				clog.Error("%s constraint.SetTimestamp err: %v, req: %v,%v", fun, err, cid, wid)
				break
			}
		case _, ok := <-channelFromMaster:
			if !ok {
				signal = api.StateNone
				return
			}
		case <-workerTicker:
			mids, exhausted, err := getMembersFunc()
			if err != nil {
				errCounter++
				clog.Error("%s getMembersFunc err: %v, req: %v,%v", fun, err, cid, wid)
				break
			}

			// exhausted means get member empty
			if exhausted {
				clog.Error("%s member tables exhausted, req: %v,%v", fun, cid, wid)
				signal = api.StateFinished
				return nil
			}

			// may empty after filtered
			if len(mids) == 0 {
				break
			}

			counter := 0
			for _, mid := range mids {
				targetModel := model.NewTarget()
				targetModel.CID = cid
				targetModel.MID = mid
				exist, err := targetModel.Upsert()
				if err != nil {
					clog.Error("%s target.Upsert err: %v, req: %v,%v,%v", fun, err, cid, wid, mid)
					errCounter++
					continue
				}
				if exist {
					clog.Warn("%s target.Upsert exist, req: %v,%v,%v", fun, cid, wid, mid)
					continue
				}

				counter++
			}

			constraintRet, err := constraintModel.IncTargetCounter(counter, -1)
			if err != nil {
				errCounter++
				clog.Error("%s constraint.IncTargetCounter err: %v, req: %v,%v", fun, err, cid, wid)
				break
			}

			if constraintRet.TargetCounter >= constraintModel.TargetTotal {
				clog.Info("%s constraint finished, req: %v,%v", fun, cid, wid)
				signal = api.StateFinished
				return nil
			}
		}
	}
}
