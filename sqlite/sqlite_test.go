package sqlite_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/Lambels/patrickarvatu.com/sqlite"
)

func TestOpenDB(t *testing.T) {
	db := MustOpenTempDB(t)
	MustCloseDB(t, db)
}

func MustOpenTempDB(t testing.TB) *sqlite.DB {
	t.Helper()

	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("tempDir: %s", err.Error())
	}
	defer os.RemoveAll(dir) // clean up temp files

	temp := filepath.Join(dir, "db")

	db := sqlite.NewDB(temp)
	if err := db.Open(); err != nil {
		t.Fatalf("open: %s", err.Error())
	}
	return db
}

func MustOpenDB(t testing.TB) *sqlite.DB {
	t.Helper()

	db := sqlite.NewDB(":memory:")
	if err := db.Open(); err != nil {
		t.Fatalf("open: %s", err.Error())
	}
	return db
}

func MustCloseDB(t testing.TB, db *sqlite.DB) {
	t.Helper()
	if err := db.Close(); err != nil {
		t.Fatalf("close: %s", err.Error())
	}
}
