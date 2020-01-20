package utils

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

// JsonTime defines a timestamp encoded as epoch seconds in JSON
type JsonTime time.Time

// MarshalJSON is used to convert the timestamp to JSON
func (t JsonTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

// UnmarshalJSON is used to convert the timestamp from JSON
func (t *JsonTime) UnmarshalJSON(s []byte) (err error) {
	r := string(s)
	r = strings.Split(r, ".")[0]
	log.Debug(r)
	q, err := strconv.ParseInt(r, 10, 64)
	if err != nil {
		return err
	}
	*(*time.Time)(t) = time.Unix(q, 0)
	return nil
}

// Unix returns t as a Unix time, the number of seconds elapsed
// since January 1, 1970 UTC. The result does not depend on the
// location associated with t.
func (t JsonTime) Unix() int64 {
	return time.Time(t).Unix()
}

// JsonTime returns the JSON time as a time.JsonTime instance in UTC
func (t JsonTime) Time() time.Time {
	return time.Time(t).UTC()
}

// String returns t as a formatted string
func (t JsonTime) String() string {
	return t.Time().String()
}
