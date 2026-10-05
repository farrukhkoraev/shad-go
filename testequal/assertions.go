
//go:build !solution

package testequal

import (
	"reflect"
)

// AssertEqual checks that expected and actual are equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are equal.
func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

    f := "%v must be equal to %v"
	args := []any {expected, actual}

	if len(msgAndArgs) > 0 {
		f = msgAndArgs[0].(string)
		args = msgAndArgs[1:]
	}

	if reflect.DeepEqual(expected, actual) {
		 return true
	} else {
		t.Errorf(f, args...)
	}
	return false
}

// AssertNotEqual checks that expected and actual are not equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are not equal.
func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	f := "%v must not be equal to %v"
	args := []any {expected, actual}

	if len(msgAndArgs) > 0 {
		f = msgAndArgs[0].(string)
		args = msgAndArgs[1:]
	}
	
	if !reflect.DeepEqual(expected, actual) {
		 return true
	} else {
		t.Errorf(f, args...)
	}
	return false
}

// RequireEqual does the same as AssertEqual but fails caller test immediately.
func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()

	f := "%v must be equal to %v"
	args := []any {expected, actual}

	if len(msgAndArgs) > 0 {
		f = msgAndArgs[0].(string)
		args = msgAndArgs[1:]
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf(f, args...)
		t.FailNow()
	}
	
}

// RequireNotEqual does the same as AssertNotEqual but fails caller test immediately.
func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()

	f := "%v must not be equal to %v"
	args := []any {expected, actual}

	if len(msgAndArgs) > 0 {
		f = msgAndArgs[0].(string)
		args = msgAndArgs[1:]
	}
	if reflect.DeepEqual(expected, actual) {
		t.Errorf(f, args...)
		t.FailNow()
	}
	
}
