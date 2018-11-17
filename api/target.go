package api

import (
	"github.com/globalsign/mgo/bson"
)

type Target struct {
	ID   bson.ObjectId    `json:"id" bson:"_id,omitempty"`
	CID  int64            `json:"cid" bson:"cid"`
	MID  int64            `json:"mid" bson:"mid"`
	PIDs map[string]State `json:"pids" bson:"pids"`
}

func NewTarget() *Target {
	return &Target{}
}
