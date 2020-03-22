package schemas

import (
	"errors"
	"fmt"
	"github.com/leocov-dev/tadpoles-backup/config"
	"strconv"
)

// custom cobra flag for concurrency value validation
type concurrencyValue int

func NewConcurrencyValue(val int, p *int) *concurrencyValue {
	*p = val
	return (*concurrencyValue)(p)
}

func (i *concurrencyValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		return err
	}

	if v > config.MaxConcurrency || v < 1 {
		return errors.New(fmt.Sprintf("value must be 1 - %d", config.MaxConcurrency))
	}

	*i = concurrencyValue(v)
	return nil
}

func (i *concurrencyValue) Type() string {
	return "int"
}

func (i *concurrencyValue) String() string { return strconv.FormatInt(int64(*i), 10) }
