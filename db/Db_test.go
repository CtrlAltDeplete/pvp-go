package db

import (
	"testing"
)

func fail(t *testing.T, msg string, expected, actual interface{}) {
	t.Fatalf("%s:\n\tExpected %v\n\tGot %v\n", msg, expected, actual)
}
