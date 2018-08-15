package project

import "github.com/simplejia/push/api"

type Project struct {
}

type Resource struct {
	PID    int64
	WID    string
	Signal api.State
}
