package itm

import (
	"fmt"
	"testing"
)

func TestRequestError(t *testing.T) {
	sut := newRequestError(fmt.Errorf("foo"))
	result := sut.Error()
	if "foo" != result {
		t.Error(unexpectedValueString("Error() result", "foo", result))
	}
}
