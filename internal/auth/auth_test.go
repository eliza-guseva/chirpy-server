package auth

import (
	"testing"
	"os"
	"os/exec"
)

func TestHashPassword_ValidOK(t *testing.T) {
	password := "ValidPassword123"

	if os.Getenv("BE_CRASHER") == "1" {
		HashPassword(password)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestHashPassword_ValidOK")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err :=cmd.Run()
	if err != nil {
		t.Errorf("Error running test TestHashPassword_ValidOK: %v", err)
	}
}

	
