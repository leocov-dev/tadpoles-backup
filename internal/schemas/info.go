package schemas

import (
	"github.com/leocov-dev/tadpoles-backup/internal/api"
	"time"
)

type Info struct {
	FirstEvent time.Time
	LastEvent  time.Time
	Dependants []string
}

func NewInfoFromParams(pr *api.ParametersResponse) *Info {
	info := &Info{
		FirstEvent: pr.FirstEventTime.Time(),
		LastEvent:  pr.LastEventTime.Time(),
	}

	for _, item := range pr.Memberships {
		for _, dep := range item.Dependants {
			info.Dependants = append(info.Dependants, dep.DisplayName)
		}
	}

	return info
}
