package proxy

import (
	"os"
	"testing"
)

func TestCheckHeader(t *testing.T) {
	f, _ := os.Open("test.dat")
	user, _ := ParseHeader(f)
	t.Log(user)
	defer f.Close()
}

func TestLoadUserFromUrl(t *testing.T) {
	loadUserFromUrl("http://localhost:8046/UsersConfig")
	t.Log("OK")
}

func TestIsValidUser(t *testing.T) {
	if isValidUser("6fto7ryo") {
		t.Log("OK")
	} else {
		t.Fail()
	}
}
