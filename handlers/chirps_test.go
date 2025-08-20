package handlers

import (
	"testing"
)

func TestCheckForProfane(t *testing.T) {
	testCases := []struct {
		chirp string
		hasProfane bool
		fixed string
	}{
		{"hello world", false, "hello world"}, // no profanity
		{"Hello sharbert", true, "Hello ****"}, // test single lowercase, preserve case
		{"hello FoRnAx", true, "hello ****"}, // test single mixedcase
		{"ForNax, fornax, fornax", true, "****, ****, ****"}, // test multiple
	}
	for _, testCase := range testCases {
		hasProfane, fixed := checkForProfane(testCase.chirp)
		if hasProfane != testCase.hasProfane {
			t.Errorf("Expected %v, got %v", testCase.hasProfane, hasProfane)
		}
		if fixed != testCase.fixed {
			t.Errorf("Expected %v, got %v", testCase.fixed, fixed)
		}
	}
}
