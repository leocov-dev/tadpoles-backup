package spinners

import (
	"fmt"
	"github.com/briandowns/spinner"
	log "github.com/sirupsen/logrus"
	"sync"
	"tadpoles-backup/config"
	"tadpoles-backup/internal/utils"
	"time"
)

// Shared Global spinnerManager
var SpinnerManager = &spinnerManager{
	lock: &sync.RWMutex{},
}

type spinnerManager struct {
	Spinners []*spinner.Spinner
	lock     *sync.RWMutex
}

func (sm *spinnerManager) AppendSpinner(s *spinner.Spinner) {
	sm.lock.Lock()
	sm.Spinners = append(sm.Spinners, s)
	sm.lock.Unlock()
}

func StartNewSpinner(title string) *spinner.Spinner {
	options := []spinner.Option{
		spinner.WithHiddenCursor(true),
		spinner.WithFinalMSG(title + " Done\n"),
	}
	s := spinner.New(config.SpinnerCharSet, config.SpinnerSpeed*time.Millisecond, options...)
	SpinnerManager.AppendSpinner(s)

	if log.GetLevel() == log.DebugLevel {
		return s
	}

	s.Prefix = fmt.Sprintf("%s ", title)
	err := s.Color("cyan", "bold") // NOTE implicit s.Start()
	if err != nil {
		utils.PrintError("Spinner startup failed: %s", err)
	}
	return s
}
