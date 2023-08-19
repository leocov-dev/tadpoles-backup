package schemas

import (
	"fmt"
	"tadpoles-backup/internal/utils"
	"time"
)

type AccountInfo struct {
	FirstEvent time.Time
	LastEvent  time.Time
	Dependants []string
}

func (i AccountInfo) PrettyPrint() {
	utils.WriteMain(
		"Time-frame",
		fmt.Sprintf(
			"%s to %s",
			i.FirstEvent.In(time.Local).Format("2006-01-02"),
			i.LastEvent.In(time.Local).Format("2006-01-02"),
		),
	)

	utils.WriteMain("Children", "")
	for i, dep := range i.Dependants {
		i += 1
		utils.WriteSub(fmt.Sprintf("%d", i), dep)
	}
}
