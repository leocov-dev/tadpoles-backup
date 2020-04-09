package utils

import (
	"math"
	"strconv"
	"time"
)

// JsonTime defines a timestamp encoded as epoch seconds in JSON
type JsonTime time.Time

func (jt JsonTime) String() string {
	return jt.Time().String()
}

// MarshalJSON is used to convert the timestamp to JSON
func (jt JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(jt).Unix(), 10)), nil
}

// UnmarshalJSON is used to convert the timestamp from JSON
func (jt *JsonTime) UnmarshalJSON(s []byte) (err error) {
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
func (jt JsonTime) Unix() int64 {
	return time.Time(jt).Unix()
}

// JsonTime returns the JSON time as a time.JsonTime instance in UTC
func (jt JsonTime) Time() time.Time {
	return (time.Time)(jt).UTC()
}
