package utils

import (
	"encoding/json"
	"testing"
	"time"
)

func TestEpocTime_String(t *testing.T) {
	expected := "2020-03-29 14:41:54 +0000 UTC"

	jt := &EpocTime{}
	err := json.Unmarshal([]byte("1585492914.00"), jt)

	if err != nil {
		t.Error(err)
	}

	if expected != jt.String() {
		t.Errorf("'%v' did not equal expected '%v'", jt.String(), expected)
	}
}

func TestEpocTime_MarshalJSON(t *testing.T) {
	expected := "1585492914.00"

	jt := EpocTime(time.Date(2020, time.March, 29, 14, 41, 54, 0, time.UTC))

	jsonBytes, err := json.Marshal(jt)
	if err != nil {
		t.Error(err)
	}

	if expected != string(jsonBytes) {
		t.Errorf("'%v' did not equal expected '%v'", string(jsonBytes), expected)
	}
}

func TestEpocTime_UnmarshalJSON(t *testing.T) {
	expected := time.Date(2020, time.March, 29, 14, 41, 54, 0, time.UTC)

	jt := &EpocTime{}
	err := json.Unmarshal([]byte("1585492914.00"), jt)

	if err != nil {
		t.Error(err)
	}

	if !expected.Equal(*(*time.Time)(jt)) {
		t.Errorf("'%v' did not equal expected '%v'", jt, expected)
	}

	if !expected.Equal(jt.Time()) {
		t.Errorf("'%v' did not equal expected '%v'", jt, expected)
	}
}

func TestEpocTime_UnmarshalJSON_Error(t *testing.T) {
	jt := &EpocTime{}
	// try to parse a string, not literal number
	err := json.Unmarshal([]byte("\"1585492914.00\""), jt)

	if err == nil {
		t.Error("Unmarshal did not raise expected error")
	}
}

func TestEpocTime_Unix(t *testing.T) {
	expected := int64(1585492914)

	jt := EpocTime(time.Date(2020, time.March, 29, 14, 41, 54, 0, time.UTC))

	asUnix := jt.Unix()

	if expected != asUnix {
		t.Errorf("'%v' did not equal expected '%v'", asUnix, expected)
	}
}
