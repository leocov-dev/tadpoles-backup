package spinners

import (
	"fmt"
	"github.com/briandowns/spinner"
	log "github.com/sirupsen/logrus"
	"sync"
	"tadpoles-backup/config"
	"time"
)

// Shared Global spinnerManager
var SpinnerManager = &spinnerManager{
	lock: &sync.RWMutex{},
}

type spinnerManager struct {
	Spinners []*WrappedSpinner
	lock     *sync.RWMutex
}
type WrappedSpinner struct {
	activeSpinner *spinner.Spinner
	Text          string
}

func (w *WrappedSpinner) Stop() {
	if w.activeSpinner != nil {
		w.activeSpinner.Stop()
	}
}

func (w *WrappedSpinner) SetPrefix(prefix string) {
	if w.activeSpinner != nil {
		w.activeSpinner.Prefix = prefix
	}
}

func (w *WrappedSpinner) Start(colors ...string) error {
	if w.activeSpinner != nil {
		return w.activeSpinner.Color(colors...) // Implicit Start()
	} else if !config.JsonOutput {
		fmt.Println(w.Text)
	}
	return nil
}

func NewWrapper(title string) *WrappedSpinner {
	w := &WrappedSpinner{
		Text: title,
	}
	if !config.NonInteractiveMode {
		options := []spinner.Option{
			spinner.WithHiddenCursor(true),
			spinner.WithFinalMSG(title + " Done\n"),
		}
		w.activeSpinner = spinner.New(config.SpinnerCharSet, config.SpinnerSpeed*time.Millisecond, options...)
	}

	w.SetPrefix(fmt.Sprintf("%s ", title))

	return w
}

func (sm *spinnerManager) AppendSpinner(s *WrappedSpinner) {
	sm.lock.Lock()
	sm.Spinners = append(sm.Spinners, s)
	sm.lock.Unlock()
}

func StartNewSpinner(title string) *WrappedSpinner {
	s := NewWrapper(title)
	SpinnerManager.AppendSpinner(s)

	if log.GetLevel() == log.DebugLevel {
		return s
	}

	err := s.Start("cyan", "bold") // NOTE implicit Start()
	if err != nil {
		log.Error("Spinner startup failed: %s", err)
	}
	return s
}
