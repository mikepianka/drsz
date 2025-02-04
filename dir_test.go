package drsz

import (
	"os"
	"path"
	"testing"
	"time"
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

func TestName(t *testing.T) {
	dir := "abc"
	dirpath := path.Join(t.TempDir(), dir)
	if err := os.Mkdir(dirpath, 0755); err != nil {
		t.Fatalf("unexpected failure: %v", err)
	}

	d, err := NewDir(dirpath)
	if err != nil {
		t.Fatalf("unexpected failure: %v", err)
	}

	if d.Name() != dir {
		t.Errorf("got %s, want %s", d.Name(), dir)
	}
}

func TestSizeString(t *testing.T) {
	tt := []struct {
		size int64
		want string
	}{
		{0, "0 B"},
		{1000, "1.0 kB"},
		{1000 * 1000, "1.0 MB"},
		{1000 * 1000 * 1000, "1.0 GB"},
		{1000 * 1000 * 1000 * 1000, "1.0 TB"},
	}

	for _, tc := range tt {
		d := Dir{SizeBytes: tc.size}
		got := d.SizeString()
		if got != tc.want {
			t.Errorf("got %s, want %s", got, tc.want)
		}
	}
}

func TestWalkCalc(t *testing.T) {
	d, err := NewDir("./testdata/example")
	if err != nil {
		t.Fatalf("unexpected failure: %v", err)
	}

	if err := d.WalkCalc(); err != nil {
		t.Fatalf("unexpected failure: %v", err)
	}

	size := int64(72379)
	if d.SizeBytes != size {
		t.Errorf("got %d, want %d", d.SizeBytes, size)
	}

	if d.LastModified.IsZero() {
		t.Errorf("LastModified should not be zero")
	}

	loc, err := time.LoadLocation("EST")
	if err != nil {
		t.Fatalf("unexpected failure: %v", err)
	}
	mod := time.Date(2025, 02, 03, 21, 50, 57, 0, loc)

	if d.LastModified.Format("2006-01-02 15:04:05") != mod.Format("2006-01-02 15:04:05") {
		t.Errorf("got %v, want %v", d.LastModified, mod)
	}
}
