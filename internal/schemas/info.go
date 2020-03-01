package schemas

import (
	"time"
)

type Info struct {
	FirstEvent time.Time
	LastEvent  time.Time
	Dependants []string
}
