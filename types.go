package gosql

import (
	"reflect"
	"time"
)

var types map[string]map[reflect.Type]string = map[string]map[reflect.Type]string{
	"sqlite3": map[reflect.Type]string{
		reflect.TypeOf(nil):         "null",
		reflect.TypeOf(0):           "integer",
		reflect.TypeOf(int64(0)):    "integer",
		reflect.TypeOf(float64(0)):  "float",
		reflect.TypeOf(false):       "integer",
		reflect.TypeOf([]byte{}):    "blob",
		reflect.TypeOf(""):          "string",
		reflect.TypeOf(time.Time{}): "datetime",
	},
}
