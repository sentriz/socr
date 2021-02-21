package db_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

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

func TestModTime(t *testing.T) {
	db, cleanup := setup(t)
	defer cleanup()

	is := is.New(t)

	date := time.Now()
	dateA := date.Add(1 * time.Hour)
	dateB := date.Add(2 * time.Hour)
	dateC := date.Add(3 * time.Hour)
	dateD := date.Add(4 * time.Hour)

	is.NoErr(db.SetModTime("laptop", "2020.01.01.jpg", dateA))
	is.NoErr(db.SetModTime("laptop", "2020.01.02.jpg", dateB))
	is.NoErr(db.SetModTime("phone", "2020.01.03.jpg", dateC))
	is.NoErr(db.SetModTime("phone", "2020.01.04.jpg", dateD))

	r, err := db.GetModTime("laptop", "2020.01.01.jpg")
	is.NoErr(err)
	is.Equal(r.Unix(), dateA.Unix())

	r, err = db.GetModTime("laptop", "2020.01.02.jpg")
	is.NoErr(err)
	is.Equal(r.Unix(), dateB.Unix())

	r, err = db.GetModTime("phone", "2020.01.03.jpg")
	is.NoErr(err)
	is.Equal(r.Unix(), dateC.Unix())

	r, err = db.GetModTime("phone", "2020.01.04.jpg")
	is.NoErr(err)
	is.Equal(r.Unix(), dateD.Unix())

	r, err = db.GetModTime("laptop", "2020.01.03.jpg")
	is.NoErr(err)
	is.Equal(r, nil)

	r, err = db.GetModTime("laptop", "2020.01.05.jpg")
	is.NoErr(err)
	is.Equal(r, nil)
}

func TestDupe(t *testing.T) {
	db, cleanup := setup(t)
	defer cleanup()

	is := is.New(t)

	date := time.Now()
	dateA := date.Add(1 * time.Hour)
	dateB := date.Add(2 * time.Hour)

	is.NoErr(db.SetModTime("a", "a", dateA))
	is.NoErr(db.SetModTime("a", "a", dateB))

	r, err := db.GetModTime("a", "a")
	is.NoErr(err)
	is.Equal(r.Unix(), dateB.Unix())
}
