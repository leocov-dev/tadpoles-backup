package utils

import (
	"testing"
	"time"
)

func TestEpocTime_UnmarshalJSON(t *testing.T) {
	expected := time.Date(2020, time.March, 29, 14, 41, 54, 0, time.UTC)

	jt := &EpocTime{}
	err := jt.UnmarshalJSON([]byte("1585492914.00"))

	if err != nil {
		t.Error(err)
	}

	if !expected.Equal(*(*time.Time)(jt)) {
		t.Errorf("'%v' did not equal expected '%v'", jt, expected)
	}
}
