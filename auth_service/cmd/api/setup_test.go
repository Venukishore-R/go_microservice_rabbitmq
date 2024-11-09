package main

import (
	"auth_service/data"
	"os"
	"testing"
)

var testApp Auth

func TestMain(m *testing.M) {
	repo := data.NewPostgresTestRepository(nil)
	testApp.Repo = repo
	os.Exit(m.Run())
}
