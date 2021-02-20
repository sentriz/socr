package db_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/matryer/is"

	"go.senan.xyz/socr/db"
)

func setup(t *testing.T) (*db.DB, func()) {
	tempFile, err := ioutil.TempFile("", "socr-test-")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	tempFilePath := tempFile.Name()
	mockDB, err := db.New(tempFilePath)
	if err != nil {
		t.Fatalf("create db: %v", err)
	}
	return mockDB, func() {
		mockDB.Close()
		os.RemoveAll(tempFilePath)
	}
}

func TestHash(t *testing.T) {
	db, cleanup := setup(t)
	defer cleanup()

	is := is.New(t)

	is.NoErr(db.SetHash("laptop", "2020.01.01.jpg", []byte("aaa")))
	is.NoErr(db.SetHash("laptop", "2020.01.02.jpg", []byte("bbb")))
	is.NoErr(db.SetHash("phone", "2020.01.03.jpg", []byte("ccc")))
	is.NoErr(db.SetHash("phone", "2020.01.04.jpg", []byte("ddd")))

	r, err := db.GetHash("laptop", "2020.01.01.jpg")
	is.NoErr(err)
	is.Equal(r, []byte("aaa"))

	r, err = db.GetHash("laptop", "2020.01.02.jpg")
	is.NoErr(err)
	is.Equal(r, []byte("bbb"))

	r, err = db.GetHash("phone", "2020.01.03.jpg")
	is.NoErr(err)
	is.Equal(r, []byte("ccc"))

	r, err = db.GetHash("phone", "2020.01.04.jpg")
	is.NoErr(err)
	is.Equal(r, []byte("ddd"))

	r, err = db.GetHash("laptop", "2020.01.03.jpg")
	is.NoErr(err)
	is.Equal(r, nil)

	r, err = db.GetHash("laptop", "2020.01.05.jpg")
	is.NoErr(err)
	is.Equal(r, nil)
}

func TestDupe(t *testing.T) {
	db, cleanup := setup(t)
	defer cleanup()

	is := is.New(t)

	is.NoErr(db.SetHash("a", "a", []byte("a")))
	is.NoErr(db.SetHash("a", "a", []byte("b")))

	r, err := db.GetHash("a", "a")
	is.NoErr(err)
	is.Equal(r, []byte("b"))
}
