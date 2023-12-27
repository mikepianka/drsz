package drsz

import "testing"

func TestCorrectlyIdentifiesCsvPath(t *testing.T) {
	p := "some/dir/a.csv"
	if !IsCsvPath(p) {
		t.Errorf("expected %s to be identified as csv", p)
	}
}

func TestCorrectlyIdentifiesNonCsvPath(t *testing.T) {
	p := "some/dir/a.txt"
	if IsCsvPath(p) {
		t.Errorf("did not expect %s to be identified as csv", p)
	}
}
