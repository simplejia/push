package constraint

import (
	"time"

	"github.com/globalsign/mgo/bson"

	"github.com/simplejia/clog/api"
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/model"
)

func WorkerMonitor(
	workerChannels map[int64]map[string]chan *Resource,
	r *Resource,
) (err error) {
	sig, cid, wid := r.Signal, r.CID, r.WID

	switch sig {
	case api.StateFinished, api.StateFailed, api.StateNone:
		delete(workerChannels[cid], wid)
		if len(workerChannels[cid]) == 0 {
			delete(workerChannels, cid)
		}

		if sig == api.StateFinished || sig == api.StateFailed {
			curState := api.StateStarted
			constraintModel := model.NewConstraint()
			constraintModel.ID = cid
			if _, err = constraintModel.UpdateState(curState, sig); err != nil {
				return
			}
		}
	}
	return
}

func StartWorker(
	workerChannels map[int64]map[string]chan *Resource,
	channelFromWorker chan *Resource,
) (err error) {
	constraintModel := model.NewConstraint()
	constraints, err := constraintModel.ListAllByState(api.StateStarted)
	if err != nil {
		return
	}

	for _, constraintModel := range constraints {
		cid := constraintModel.ID

		if constraintModel.WorkerCounter >= constraintModel.WorkerTotal {
			continue
		}

		// if worker not exit yet
		if len(workerChannels[cid]) >= constraintModel.WorkerTotal {
			continue
		}

		for {
			constraintRet, err := constraintModel.IncWorkerCounter(1, constraintModel.WorkerTotal)
			if err != nil {
				return err
			}

			// if already started enough workers
			if constraintRet == nil {
				break
			}

			if workerChannels[cid] == nil {
				workerChannels[cid] = make(map[string]chan *Resource)
			}
			wid := bson.NewObjectId().Hex()
			channelToWorker := make(chan *Resource, 10)
			workerChannels[cid][wid] = channelToWorker

			go Worker(cid, wid, channelToWorker, channelFromWorker)

			// avoid all worker in one server
			r := time.Now().UnixNano() % 1000
			time.Sleep(time.Millisecond * time.Duration(r))
		}
	}

	return
}

func StopWorker(
	workerChannels map[int64]map[string]chan *Resource,
) (err error) {
	for cid := range workerChannels {
		constraintModel := model.NewConstraint()
		constraintModel.ID = cid
		constraintModel, err := constraintModel.Get()
		if err != nil {
			return err
		}

		switch constraintModel.State {
		case api.StateStopped, api.StateFinished, api.StateFailed, api.StateReboot, api.StatePause:
			for wid, channel := range workerChannels[cid] {
				if channel == nil {
					continue
				}

				// just close, delete when get notified when worker exit really
				close(channel)
				workerChannels[cid][wid] = nil
			}
		}
	}

	return
}

func RebootWorker() (err error) {
	if constraints, err := model.NewConstraint().ListAllByState(api.StateStarted); err != nil {
		return err
	} else {
		for _, constraintModel := range constraints {
			if _, hasDead := constraintModel.CheckAllWorkers(); hasDead {
				// if some worker is dead, change state to reboot
				if _, err := constraintModel.UpdateState(api.StateStarted, api.StateReboot); err != nil {
					return err
				}
			}
		}
	}

	if constraints, err := model.NewConstraint().ListAllByState(api.StateReboot); err != nil {
		return err
	} else {
		for _, constraintModel := range constraints {
			if hasAlive, _ := constraintModel.CheckAllWorkers(); !hasAlive {
				// if all worker dead, change state to started
				if _, err := constraintModel.Reboot(); err != nil {
					return err
				}
			}
		}
	}

	return
}

func Master() {
	fun := "constraint.Master"

	channelFromWorker := make(chan *Resource, 100)
	workerChannels := make(map[int64]map[string]chan *Resource) // cid->wid->Resource
	ticker := time.Tick(time.Second)

	for {
		select {
		case <-ticker:
			if err := StartWorker(workerChannels, channelFromWorker); err != nil {
				clog.Error("%s StartWorker err: %v", fun, err)
				break
			}
			if err := StopWorker(workerChannels); err != nil {
				clog.Error("%s StopWorker err: %v", fun, err)
				break
			}
			if err := RebootWorker(); err != nil {
				clog.Error("%s RebootWorker err: %v", fun, err)
				break
			}
		case r := <-channelFromWorker:
			if err := WorkerMonitor(workerChannels, r); err != nil {
				clog.Error("%s WorkerMonitor err: %v", fun, err)
				break
			}
		}
	}
}

func init() {
	go Master()
}
