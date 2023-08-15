package schemas

import (
	"fmt"
	"sort"
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

func NewInfoFromEvents(events api.Events) *Info {
	info := &Info{
		FirstEvent: events[0].EventTime.Time(),
		LastEvent:  events[len(events)-1].EventTime.Time(),
	}

	childSet := make(map[string]bool)

	for _, event := range events {
		if event.ChildName == "" {
			continue
		}
		childSet[event.ChildName] = true
	}

	for name, _ := range childSet {
		info.Dependants = append(info.Dependants, name)
	}

	sort.Strings(info.Dependants)

	return info
}
