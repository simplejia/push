package constraint

import "github.com/simplejia/push/api"

//Constraint struct in service package
type Constraint struct {
}

type Resource struct {
	CID    int64
	WID    string
	Signal api.State
}
