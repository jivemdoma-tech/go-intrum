package gointrum

import "testing"

func assertNoErr(t *testing.T, err error) {
	if err != nil {
		t.Errorf("unexpected error: %v", err.Error())
	}
}

func assertErr(t *testing.T, err error) {
	if err == nil {
		t.Errorf("error expected")
	}
}
