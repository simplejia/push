package project

import (
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/simplejia/clog"
	"github.com/simplejia/push/api"
	"github.com/simplejia/push/model"
)

func WorkerMonitor(
	workerChannels map[int64]map[string]chan *Resource,
	r *Resource,
) (err error) {
	sig, pid, wid := r.Signal, r.PID, r.WID

	switch sig {
	case api.StateFinished, api.StateFailed, api.StateNone:
		delete(workerChannels[pid], wid)
		if len(workerChannels[pid]) == 0 {
			delete(workerChannels, pid)
		}

		if sig == api.StateFinished || sig == api.StateFailed {
			curState := api.StateStarted
			projectModel := model.NewProject()
			projectModel.ID = pid
			if _, err = projectModel.UpdateState(curState, sig); err != nil {
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
	projectModel := model.NewProject()
	projects, err := projectModel.ListAllByState(api.StateStarted)
	if err != nil {
		return
	}

	for _, projectModel := range projects {
		pid := projectModel.ID

		if projectModel.WorkerCounter >= projectModel.WorkerTotal {
			continue
		}

		// if worker not exit yet
		if len(workerChannels[pid]) >= projectModel.WorkerTotal {
			continue
		}

		for {
			projectRet, err := projectModel.IncWorkerCounter(1, projectModel.WorkerTotal)
			if err != nil {
				return err
			}

			// if already started enough workers
			if projectRet == nil {
				break
			}

			//insert channel to workerChannels
			if workerChannels[pid] == nil {
				workerChannels[pid] = make(map[string]chan *Resource)
			}
			wid := bson.NewObjectId().Hex()
			channelToWorker := make(chan *Resource, 10)
			workerChannels[pid][wid] = channelToWorker

			go Worker(pid, wid, channelToWorker, channelFromWorker)

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
	for pid := range workerChannels {
		projectModel := model.NewProject()
		projectModel.ID = pid
		projectModel, err := projectModel.Get()
		if err != nil {
			return err
		}

		switch projectModel.State {
		case api.StateStopped, api.StateFinished, api.StateFailed, api.StateReboot, api.StatePause:
			for wid, channel := range workerChannels[pid] {
				if channel == nil {
					continue
				}

				// just close, delete when get notified when worker exit really
				close(channel)
				workerChannels[pid][wid] = nil
			}
		}
	}

	return
}

func RebootWorker() (err error) {
	if projects, err := model.NewProject().ListAllByState(api.StateStarted); err != nil {
		return err
	} else {
		for _, projectModel := range projects {
			if _, hasDead := projectModel.CheckAllWorkers(); hasDead {
				// if some worker is dead, change state to reboot
				if _, err := projectModel.UpdateState(api.StateStarted, api.StateReboot); err != nil {
					return err
				}
			}
		}
	}

	if projects, err := model.NewProject().ListAllByState(api.StateReboot); err != nil {
		return err
	} else {
		for _, projectModel := range projects {
			if hasAlive, _ := projectModel.CheckAllWorkers(); !hasAlive {
				// if all worker dead, change state to started
				if _, err := projectModel.Reboot(); err != nil {
					return err
				}
			}
		}
	}

	return
}

func Master() {
	fun := "project.Master"

	channelFromWorker := make(chan *Resource, 100)
	workerChannels := make(map[int64]map[string]chan *Resource) // pid->wid->Resource
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
