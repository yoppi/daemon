package daemon

import (
	"reflect"
	"testing"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func Test_Daemonize(t *testing.T) {
	d := Daemonize("test.log")
	expect(t, d.Logfile, "test.log")
}
