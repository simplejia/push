package project

import (
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/simplejia/push/api"
)

func (project *Project) UpdateState(curState, newState api.State) (exist bool, err error) {
	c := project.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"_id": project.ID,
	}

	if curState != api.StateNone {
		q["state"] = curState
	}

	change := mgo.Change{
		Update: bson.M{
			"$set": bson.M{
				"state": newState,
			},
		},
		ReturnNew: true,
	}

	info, err := c.Find(q).Apply(change, nil)
	if err != nil {
		if err != mgo.ErrNotFound {
			return
		}

		err = nil
		return
	}
	if info.Matched == 1 {
		exist = true
	}

	return
}

func (project *Project) SetTimestamp(wid string, timestamp int64) (err error) {
	c := project.GetC()
	defer c.Database.Session.Close()

	Update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("timestamps.%s", wid): timestamp,
		},
	}

	err = c.UpdateId(project.ID, Update)
	if err != nil {
		return
	}
	return

}

func (project *Project) UnsetTimestamp(wid string) (err error) {
	c := project.GetC()
	defer c.Database.Session.Close()

	unsetField := fmt.Sprintf("timestamps.%s", wid)
	if wid == "" {
		unsetField = "timestamps"
	}

	Update := bson.M{
		"$unset": bson.M{
			unsetField: nil,
		},
	}

	err = c.UpdateId(project.ID, Update)
	if err != nil {
		return
	}

	return
}

func (project *Project) Reboot() (exist bool, err error) {
	c := project.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"_id":   project.ID,
		"state": api.StateReboot,
	}

	change := mgo.Change{
		Update: bson.M{
			"$unset": bson.M{
				"timestamps": nil,
			},
			"$set": bson.M{
				"state":          api.StateStarted,
				"worker_counter": 0,
			},
		},
	}

	info, err := c.Find(q).Apply(change, nil)
	if err != nil {
		if err != mgo.ErrNotFound {
			return
		}

		err = nil
		return
	}
	if info.Matched == 1 {
		exist = true
	}

	return
}

func (project *Project) Start() (exist bool, err error) {
	c := project.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"_id": project.ID,
		"state": bson.M{
			"$ne": api.StateStarted,
		},
	}

	change := mgo.Change{
		Update: bson.M{
			"$unset": bson.M{
				"timestamps": nil,
			},
			"$set": bson.M{
				"state":          api.StateStarted,
				"worker_counter": 0,
				"target_counter": 0,
				"started_time":   time.Now().Unix(),
			},
		},
	}

	info, err := c.Find(q).Apply(change, nil)
	if err != nil {
		if err != mgo.ErrNotFound {
			return
		}

		err = nil
		return
	}
	if info.Matched == 1 {
		exist = true
	}

	return
}

func (project *Project) Proceed() (exist bool, err error) {
	c := project.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"_id":   project.ID,
		"state": api.StatePause,
	}

	change := mgo.Change{
		Update: bson.M{
			"$unset": bson.M{
				"timestamps": nil,
			},
			"$set": bson.M{
				"state":          api.StateStarted,
				"worker_counter": 0,
			},
		},
	}

	info, err := c.Find(q).Apply(change, nil)
	if err != nil {
		if err != mgo.ErrNotFound {
			return
		}

		err = nil
		return
	}
	if info.Matched == 1 {
		exist = true
	}

	return
}
