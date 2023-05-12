package storage_test

import (
	"os"
	"testing"
)

var (
	dbUrl string
)

func TestMain(m *testing.M) {
	dbUrl = os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "host=localhost dbname=PhoneTracker_test sslmode=disable user=postgres password=hnxJpVsk3r"
	}

	os.Exit(m.Run())
}
