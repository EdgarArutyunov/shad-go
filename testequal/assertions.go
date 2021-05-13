// +build !solution

package testequal

import (
	"fmt"
	"reflect"
)

// AssertEqual checks that expected and actual are equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are equal.

func eq(a, b interface{}) (bool, error) {
	if a == nil && b == nil {
		return a == b, nil
	}

	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false, nil
	}

	switch a.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		string:

		return a == b, nil

	case []int:
		arA := a.([]int)
		arB := b.([]int)

		if arA == nil || arB == nil {
			return arA == nil && arB == nil, nil
		}

		if len(arA) != len(arB) {
			return false, nil
		}

		for i := range arA {
			if arA[i] != arB[i] {
				return false, nil
			}
		}

		return true, nil

	case []byte:
		arA := a.([]byte)
		arB := b.([]byte)

		if arA == nil || arB == nil {
			return arA == nil && arB == nil, nil
		}

		if len(arA) != len(arB) {
			return false, nil
		}

		for i := range arA {
			if arA[i] != arB[i] {
				return false, nil
			}
		}

		return true, nil

	case map[string]string:
		mA := a.(map[string]string)
		mB := b.(map[string]string)

		if mA == nil || mB == nil {
			return mA == nil && mB == nil, nil
		}

		if len(mA) != len(mB) {
			return false, nil
		}

		for key, val := range mA {
			if vb, ok := mB[key]; !ok || vb != val {
				return false, nil
			}
		}

		return true, nil
	default:
		return false, fmt.Errorf("unsuported type")
	}
}

func AssertEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	if equals, err := eq(expected, actual); !equals || err != nil {
		if len(msgAndArgs) == 0 {
			return false
		}

		format, ok := msgAndArgs[0].(string)
		if !ok {
			panic("not correct format string")
		}

		t.Errorf(format, msgAndArgs[1:]...)
		return false
	}
	return true
}

// AssertNotEqual checks that expected and actual are not equal.
//
// Marks caller function as having failed but continues execution.
//
// Returns true iff arguments are not equal.
func AssertNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
	t.Helper()

	if equals, err := eq(expected, actual); equals || err != nil {
		if len(msgAndArgs) == 0 {
			return false
		}

		format, ok := msgAndArgs[0].(string)
		if !ok {
			panic("not correct format string")
		}

		t.Errorf(format, msgAndArgs[1:]...)
		return false
	}
	return true
}

// RequireEqual does the same as AssertEqual but fails caller test immediately.
func RequireEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()

	if equals, err := eq(expected, actual); !equals || err != nil {
		if len(msgAndArgs) == 0 {
			t.Errorf("")
			t.FailNow()
		}

		format, ok := msgAndArgs[0].(string)
		if !ok {
			panic("not correct format string")
		}

		t.Errorf(format, msgAndArgs[1:]...)
		t.FailNow()
	}
}

// RequireNotEqual does the same as AssertNotEqual but fails caller test immediately.
func RequireNotEqual(t T, expected, actual interface{}, msgAndArgs ...interface{}) {
	t.Helper()

	if equals, err := eq(expected, actual); equals && err == nil {

		if len(msgAndArgs) == 0 {
			t.Errorf("")
			t.FailNow()
		}

		format, ok := msgAndArgs[0].(string)
		if !ok {
			panic("not correct format string")
		}

		t.Errorf(format, msgAndArgs[1:]...)
		t.FailNow()
	}
}
