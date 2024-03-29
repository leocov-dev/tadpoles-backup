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
	Spinners []WrappedSpinner
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
	if w.activeSpinner != nil { // nil would be non-interactive
		return w.activeSpinner.Color(colors...) // Implicit Start()
	} else if config.IsHumanReadable() { // not json but still non-interactive
		fmt.Println(w.Text)
	}
	return nil
}

func NewWrapper(title string) *WrappedSpinner {
	w := &WrappedSpinner{
		Text: title,
	}
	if config.IsInteractive() {
		options := []spinner.Option{
			spinner.WithHiddenCursor(true),
		}
		w.activeSpinner = spinner.New(config.SpinnerCharSet, config.SpinnerSpeed*time.Millisecond, options...)
	}

	w.SetPrefix(fmt.Sprintf("%s ", title))

	return w
}

func (sm *spinnerManager) AppendSpinner(s *WrappedSpinner) {
	sm.lock.Lock()
	sm.Spinners = append(sm.Spinners, *s)
	sm.lock.Unlock()
}

func (sm *spinnerManager) StopAll() {
	sm.lock.Lock()
	for _, ws := range sm.Spinners {
		ws.Stop()
	}
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
		log.Errorf("Spinner startup failed: %s", err)
	}
	return s
}
