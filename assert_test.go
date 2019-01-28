// Assertion helpers of golang unit tests
package main

import (
	"fmt"
	"runtime/debug"
	"testing"
)

func assertTrue(t *testing.T, message string, check bool) {
	if !check {
		debug.PrintStack()
		t.Fatal(message)
	}
}

func assertFalse(t *testing.T, message string, check bool) {
	if check {
		debug.PrintStack()
		t.Fatal(message)
	}
}

func assertExpectNoErr(t *testing.T, message string, err error) {
	if err != nil {
		debug.PrintStack()
		t.Fatalf("%s : %s", message, err)
	}
}

func assertExpectErr(t *testing.T, message string, err error) {
	if err == nil {
		debug.PrintStack()
		t.Fatal(message)
	}
}

func assertEqualsInt(t *testing.T, message string, expected int, actual int) {
	assertTrue(t, fmt.Sprintf("%s\nExpected: %d, Actual: %d", message, expected, actual), expected == actual)
}

func assertEqualsStr(t *testing.T, message string, expected string, actual string) {
	assertTrue(t, fmt.Sprintf("%s\nExpected: %s, Actual: %s", message, expected, actual), expected == actual)
}

func assertEqualsBool(t *testing.T, message string, expected bool, actual bool) {
	assertTrue(t, fmt.Sprintf("%s\nExpected: %t, Actual: %t", message, expected, actual), expected == actual)
}

func assertEqualsSlice(t *testing.T, message string, expected []uint32, actual []uint32) {
	assertEqualsInt(t, fmt.Sprintf("%s\nSize missmatch", message), len(expected), len(actual))
	for index, expvalue := range expected {
		actvalue := actual[index]
		assertTrue(t, fmt.Sprintf("%s\nIndex %d - Expected: %d, Actual: %d", message, index, expvalue,
			actvalue), expvalue == actvalue)
	}
}
