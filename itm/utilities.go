package itm

import (
	"fmt"
)

func unexpectedValueString(label string, expected interface{}, got interface{}) string {
	return fmt.Sprintf("Unexpected value [%s]\nExpected: %#v\n     Got: %#v", label, expected, got)
}
