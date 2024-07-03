package main_test

import (
	"testing"

	db_test "github.com/RyoheiTomiyama/phraze-api/test/db"
)

func TestMain(m *testing.M) {
	_, err := db_test.SetupDB()
	if err != nil {
		panic(err)
	}

	m.Run()

}
