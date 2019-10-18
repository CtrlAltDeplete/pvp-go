package db

import (
	"testing"
)

func fail(t *testing.T, functionName string, expected, actual interface{}) {
	t.Fatalf("%s failed:\n\tExpected %v\n\tGot %v\n", functionName, expected, actual)
}
