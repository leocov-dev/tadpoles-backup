package login

import "time"

type Login interface {
	NeedsLogin() bool
	DoLogin(user string, password string) (*time.Time, error)
}
