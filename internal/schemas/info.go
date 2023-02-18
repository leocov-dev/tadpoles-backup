package schemas

import (
	"fmt"
	"tadpoles-backup/internal/api"
	"tadpoles-backup/internal/utils"
	"time"
)

type Info struct {
	FirstEvent time.Time `json:"firstEvent"`
	LastEvent  time.Time `json:"lastEvent"`
	Dependants []string  `json:"dependants"`
}

func (i Info) prettyFormatTimeFrame() string {
	return fmt.Sprintf(
		"%s to %s",
		i.FirstEvent.In(time.Local).Format("2006-01-02"),
		i.LastEvent.In(time.Local).Format("2006-01-02"),
	)
}

func (i Info) PrettyPrint() {
	utils.WriteMain("Time-frame", i.prettyFormatTimeFrame())

	utils.WriteMain("Children", "")
	for i, dep := range i.Dependants {
		i += 1
		utils.WriteSub(fmt.Sprintf("%d", i), dep)
	}
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
