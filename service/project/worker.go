package project

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/simplejia/clog"
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/model"
)

func Push(msg json.RawMessage, mids []int64) (failMids map[int64]bool, err error) {
	// TODO: begin push
	// ...
	return
}

func Worker(pid int64, wid string, channelFromMaster, channelToMaster chan *Resource) (err error) {
	fun := "project.Worker"
	clog.Info("%s worker started, req: %v,%v", fun, pid, wid)

	var signal api.State

	defer func() {
		channelToMaster <- &Resource{
			PID:    pid,
			WID:    wid,
			Signal: signal,
		}

		if err != nil {
			clog.Error("%s err: %v, req: %v,%v", fun, err, pid, wid)
		}
	}()

	heartbeatTicker := time.Tick(api.KeepaliveTimeThreshold)
	errCounter, errThreshold := 0, api.ErrorThreshold

	projectModel := model.NewProject()
	projectModel.ID = pid
	projectModel, err = projectModel.Get()
	if err != nil {
		signal = api.StateFailed
		return
	}

	if projectModel == nil {
		err = errors.New("project.Get empty")
		signal = api.StateFailed
		return
	}

	qps, targetNum := projectModel.QPS, 50
	if qps <= 0 {
		qps = 10000
	}

	workerTicker := time.Tick(time.Duration(int(time.Second) * targetNum * projectModel.WorkerTotal / qps))

	for {
		if errCounter >= errThreshold {
			err = errors.New("failed, err threshold exceeded")
			signal = api.StateFailed
			return
		}

		select {
		case <-heartbeatTicker: // send heartbeat timestamp.
			if err := projectModel.SetTimestamp(wid, time.Now().Unix()); err != nil {
				errCounter++
				clog.Error("%s project.SetTimestamp err: %v, req: %v,%v", fun, err, pid, wid)
			}
		case _, ok := <-channelFromMaster:
			if !ok {
				signal = api.StateNone
				return
			}
		case <-workerTicker:
			targetModel := model.NewTarget()
			targetModel.CID = projectModel.CID
			targets, err := targetModel.ListAndSet(pid, targetNum)
			if err != nil {
				errCounter++
				clog.Error("%s target.ListAdSet err: %v, req: %v,%v", fun, err, pid, wid)
				break
			}

			if len(targets) == 0 {
				clog.Info("%s project finished(no target), req: %v,%v", fun, pid, wid)
				signal = api.StateFinished
				return nil
			}

			mids := []int64{}
			for _, targetModel := range targets {
				mids = append(mids, targetModel.MID)
			}

			msg := projectModel.Msg

			failMids, err := Push(msg, mids)
			if err != nil {
				clog.Error("%s Push err: %v, req: %v,%v,%s,%v", fun, err, pid, wid, msg, mids)
				errCounter++
				break
			}

			for _, targetModel := range targets {
				state := api.StateFinished
				if failMids[targetModel.MID] {
					state = api.StateFailed
				}

				pids := map[string]api.State{}
				pids[strconv.FormatInt(pid, 10)] = state

				targetModel.PIDs = pids // reuse pids

				if err := targetModel.SetPidState(); err != nil {
					clog.Error("%s target.SetPidState err: %v, req: %v,%v,%v", fun, err, pid, wid, targetModel.MID)
					errCounter++
					continue
				}
			}

			counter := len(mids) - len(failMids)

			//increase the pushed counter
			projectRet, err := projectModel.IncTargetCounter(counter, -1)
			if err != nil {
				clog.Error("%s project.IncTargetCounter err: %v, req: %v,%v", fun, err, pid, wid)
				errCounter++
				break
			}

			if projectRet.TargetCounter >= projectModel.TargetTotal {
				clog.Info("%s project finished(count exceeded), req: %v,%v", fun, pid, wid)
				signal = api.StateFinished
				return nil
			}
		}
	}
}
