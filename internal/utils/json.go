package utils

import (
	"math"
	"strconv"
	"time"
)

// EpocTime defines a timestamp encoded as epoch seconds in JSON
type EpocTime time.Time

func (jt EpocTime) String() string {
	return jt.Time().String()
}

// MarshalJSON is used to convert the timestamp to JSON
func (jt EpocTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(time.Time(jt).Unix()), 'f', 2, 64)), nil
}

// UnmarshalJSON is used to convert the timestamp from JSON
func (jt *EpocTime) UnmarshalJSON(s []byte) (err error) {
	r := string(s)

	f, err := strconv.ParseFloat(r, 64)
	if err != nil {
		return err
	}

	sec, dec := math.Modf(f)
	*(*time.Time)(jt) = time.Unix(int64(sec), int64(dec*(1e9)))

	return nil
}

// Unix returns t as a Unix time, the number of seconds elapsed
// since January 1, 1970 UTC. The result does not depend on the
// location associated with t.
func (jt EpocTime) Unix() int64 {
	return time.Time(jt).Unix()
}

// Time returns the JSON time as a time.EpocTime instance in UTC
func (jt EpocTime) Time() time.Time {
	return (time.Time)(jt).UTC()
}
