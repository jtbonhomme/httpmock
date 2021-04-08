package httpmock

import (
	"bytes"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// Assert fails the test if the condition is false.
func Assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		b := bytes.NewBufferString("\t" + msg + "\n")
		fmt.Fprintln(b, v...)
		print(b)
		tb.FailNow()
	}
}

// OK fails the test if an err is not nil.
func OK(tb testing.TB, err error) {
	if err != nil {
		print(bytes.NewBufferString(
			fmt.Sprintf("\tUnexpected error: %v", err)))
		tb.FailNow()
	}
}

// NotNil fails the test if anything is nil.
func NotNil(tb testing.TB, anything interface{}) {
	if isNil(anything) {
		print(bytes.NewBufferString("\tExpected non-nil value"))
		tb.FailNow()
	}
}

// Nil fails the test if something is NOT nil.
func Nil(tb testing.TB, something interface{}) {
	if !isNil(something) {
		print(bytes.NewBufferString(
			fmt.Sprintf("\tExpected value to be nil\n\n\tgot: %#v", something)))
		tb.FailNow()
	}
}

// Equals fails the test if exp is not equal to act.
func Equals(tb testing.TB, exp, act interface{}) {
	if b, ok := equals(exp, act); !ok {
		print(b)
		tb.FailNow()
	}
}

func equals(exp, act interface{}) (b *bytes.Buffer, ok bool) {
	b = new(bytes.Buffer)
	fmt.Fprintf(b, "\texp: %s\n\n\tgot: %s", stringer(exp), stringer(act))
	return b, reflect.DeepEqual(exp, act)
}

// Includes fails if expected string is NOT included in the actual string
func Includes(tb testing.TB, exp string, act ...string) {
	for _, a := range act {
		if strings.Contains(a, exp) {
			return
		}
	}

	print(bytes.NewBufferString(
		fmt.Sprintf("\tExpected to include: %s\n\n\tgot: %s", exp, act)))
	tb.FailNow()
}

// NotIncludes fails if expected string is included in the actual string
func NotIncludes(tb testing.TB, exp string, act ...string) {
	for _, a := range act {
		if strings.Contains(a, exp) {
			print(bytes.NewBufferString(
				fmt.Sprintf("\tNOT expected to include: %#v\n\n\tgot: %#v", exp, act)))
			tb.FailNow()
		}
	}
}

// IncludesI fails if expected string is NOT included in the actuall string (ignore case)
func IncludesI(tb testing.TB, exp string, act ...string) {
	for _, a := range act {
		if strings.Contains(strings.ToLower(a), strings.ToLower(exp)) {
			return
		}
	}

	print(bytes.NewBufferString(
		fmt.Sprintf("\tExpected to include: %s\n\n\tgot: %s", exp, act)))
	tb.FailNow()
}

// IncludesSlice fails if all of expected items is NOT included in the actual slice
func IncludesSlice(tb testing.TB, exp, act interface{}) {
	if reflect.ValueOf(exp).Kind() != reflect.Slice {
		panic("IncludesSlice requires a expected slice")
	}

	if reflect.ValueOf(act).Kind() != reflect.Slice {
		panic("IncludesSlice requires a actual slice")
	}

	expSlice := reflect.ValueOf(exp)
	actSlice := reflect.ValueOf(act)

	expLen := expSlice.Len()
	actLen := actSlice.Len()

	if expLen <= actLen {
		var score int
		for idxA := 0; idxA < actLen; idxA++ {
			for idxE := 0; idxE < expLen; idxE++ {
				if reflect.DeepEqual(expSlice.Index(idxE).Interface(), actSlice.Index(idxA).Interface()) {
					score++
				}
			}
		}
		if score == expLen {
			return
		}
	}

	print(bytes.NewBufferString(
		fmt.Sprintf("\tExpected to all items to be included: %+v\n\n\tIn: %+v", exp, act)))
	tb.FailNow()
}

// IncludesMap fails if all of expected map entries are NOT included in the actuall map
func IncludesMap(tb testing.TB, exp, act interface{}) {
	if ok := includesMap(exp, act); !ok {
		tb.FailNow()
	}
}

func includesMap(exp, act interface{}) (ok bool) {
	if reflect.ValueOf(exp).Kind() != reflect.Map {
		panic("IncludesMap requires a expected map")
	}

	if reflect.ValueOf(act).Kind() != reflect.Map {
		panic("IncludesMap requires a actual map")
	}

	expMap := reflect.ValueOf(exp)
	actMap := reflect.ValueOf(act)

	expLen := len(expMap.MapKeys())
	actLen := len(actMap.MapKeys())

	if expLen <= actLen {
		var score int
		for _, actKey := range actMap.MapKeys() {
			for _, expKey := range expMap.MapKeys() {
				if reflect.DeepEqual(expKey.Interface(), actKey.Interface()) &&
					reflect.DeepEqual(expMap.MapIndex(expKey).Interface(), actMap.MapIndex(actKey).Interface()) {
					score++
				}
			}
		}

		if score == expLen {
			return true
		}
	}

	return false
}

// Zero fails the test if anything is NOT nil.
func Zero(tb testing.TB, anything interface{}) {
	if !isZero(anything) {
		print(bytes.NewBufferString("\tExpected zero value"))
		tb.FailNow()
	}
}

// NotZero fails the test if anything is NOT nil.
func NotZero(tb testing.TB, anything interface{}) {
	if isZero(anything) {
		print(bytes.NewBufferString("\tExpected non-zero value"))
		tb.FailNow()
	}
}

func isZero(anything interface{}) bool {
	refZero := reflect.Zero(reflect.ValueOf(anything).Type())
	return reflect.DeepEqual(refZero.Interface(), anything)
}

func isNil(anything interface{}) bool {
	return reflect.DeepEqual(reflect.ValueOf(nil), reflect.ValueOf(anything)) ||
		reflect.ValueOf(anything).IsNil()
}

func print(b *bytes.Buffer) {
	_, file, line, _ := runtime.Caller(2)
	fmt.Printf("\033[31m%s:%d:\n\n%s\033[39m\n\n",
		filepath.Base(file), line, b.String())
}

func stringer(a interface{}) string {
	switch s := a.(type) {
	case string:
		return s
	case []byte:
		return string(s)
	default:
		return fmt.Sprintf("%#v", s)
	}
}
