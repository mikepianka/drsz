package drsz

import (
	"testing"
)

func TestNewDir(t *testing.T) {
	d, err := NewDir(".")
	if err != nil {
		t.Errorf("NewDir() failed: %v", err)
	}

	if d.AbsPath == "" {
		t.Errorf("NewDir() failed: AbsPath should not be empty")
	}

	if d.SizeBytes != 0 {
		t.Errorf("NewDir() failed: SizeBytes expected to be zero")
	}

	if !d.LastModified.IsZero() {
		t.Errorf("NewDir() failed: LastModified expected to be zero")
	}
}
