package target

import (
	"fmt"

	"github.com/simplejia/push/api"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (target *Target) ListAndSet(pid int64, limit int) (targets []*Target, err error) {
	c := target.GetC()
	defer c.Database.Session.Close()

	mkey := fmt.Sprintf("pids.%d", pid)

	q := bson.M{
		"cid": target.CID,
		mkey: bson.M{
			"$exists": false,
		},
	}

	change := mgo.Change{
		Update: bson.M{
			"$set": bson.M{
				mkey: api.StateStarted,
			},
		},
		ReturnNew: true,
	}

	for i := 0; i < limit; i++ {
		var updatedTarget *Target
		_, err = c.Find(q).Apply(change, &updatedTarget)
		if err != nil {
			if err != mgo.ErrNotFound {
				return
			}
			err = nil
			return
		}
		targets = append(targets, updatedTarget)
	}

	return
}

func (target *Target) ListWithPidState(pid int64, state api.State, offset, limit int) (targets []*Target, err error) {
	c := target.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"cid": target.CID,
		fmt.Sprintf("pids.%d", pid): state,
	}

	if err = c.Find(q).Sort("-_id").Skip(offset).Limit(limit).All(&targets); err != nil {
		return
	}

	return
}

func (target *Target) List(offset, limit int) (targets []*Target, err error) {
	c := target.GetC()
	defer c.Database.Session.Close()

	err = c.Find(nil).Sort("-_id").Skip(offset).Limit(limit).All(&targets)
	if err != nil {
		return
	}

	return
}

func (target *Target) ListByID(limit int) (targets []*Target, err error) {
	c := target.GetC()
	defer c.Database.Session.Close()

	q := bson.M{}

	if id := target.ID; id != "" {
		q["_id"] = bson.M{
			"$gt": id,
		}
	}
	err = c.Find(q).Sort("_id").Limit(limit).All(&targets)
	if err != nil {
		return
	}

	return
}

func (target *Target) ListByCid(cid int64, offset, limit int) (targets []*Target, err error) {
	c := target.GetC()
	defer c.Database.Session.Close()

	q := bson.M{
		"cid": cid,
	}
	err = c.Find(q).Sort("-_id").Skip(offset).Limit(limit).All(&targets)
	if err != nil {
		return
	}

	return
}
