package progress

import (
	"fmt"
	"github.com/gosuri/uiprogress"
	"tadpoles-backup/config"
)

type BarWrapper struct {
	bar *uiprogress.Bar
}

func (bw *BarWrapper) Stop() {
	if bw.bar != nil {
		uiprogress.Stop()
	}
}

func (bw *BarWrapper) Increment() {
	if bw.bar != nil {
		bw.bar.Incr()
	}
}

func StartNewProgressBar(steps int, heading string) *BarWrapper {
	wrapper := &BarWrapper{}
	if !config.NonInteractiveMode {
		uiprogress.Start()
		wrapper.bar = uiprogress.AddBar(steps).
			AppendCompleted().
			PrependElapsed().
			PrependFunc(func(b *uiprogress.Bar) string {
				return fmt.Sprintf("%s (%d/%d)", heading, b.Current(), steps)
			})
	}
	return wrapper
}
