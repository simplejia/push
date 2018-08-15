package constraint

import (
	"fmt"
	"time"

	"github.com/simplejia/push/api"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (constraint *Constraint) UpdateState(curState, newState api.State) (exist bool, err error) {
	c := constraint.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"_id": constraint.ID,
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

func (constraint *Constraint) SetTimestamp(wid string, timestamp int64) (err error) {
	c := constraint.GetC()
	defer c.Database.Session.Close()

	Update := bson.M{
		"$set": bson.M{
			fmt.Sprintf("timestamps.%s", wid): timestamp,
		},
	}

	err = c.UpdateId(constraint.ID, Update)
	if err != nil {
		return
	}

	return
}

func (constraint *Constraint) UnsetTimestamp(wid string) (err error) {
	c := constraint.GetC()
	defer c.Database.Session.Close()

	unsetField := fmt.Sprintf("timestamps.%s", wid)
	if wid == "" {
		unsetField = "timestamps"
	}

	update := bson.M{
		"$unset": bson.M{
			unsetField: nil,
		},
	}

	err = c.UpdateId(constraint.ID, update)
	if err != nil {
		return
	}

	return
}

func (constraint *Constraint) Reboot() (exist bool, err error) {
	c := constraint.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"_id":   constraint.ID,
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

func (constraint *Constraint) Start() (exist bool, err error) {
	c := constraint.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"_id": constraint.ID,
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

func (constraint *Constraint) Proceed() (exist bool, err error) {
	c := constraint.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"_id":   constraint.ID,
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
