package errors

import (
	"fmt"
	"testing"
)

func TestErr(t *testing.T) {
	err := New("test")
	err = Abc(err)
	err = ddd(err)
	trace := GetStackTrace(err)
	fmt.Println(trace)
}
func Abc(err error) error {
	return Wrap(err, "abc")
}
func ddd(err error) error {
	return err
}
