package models

import (
	"reflect"
	"testing"
)

func TestTypes(t *testing.T) {
	for _, v := range Types {
		reflectType := reflect.TypeOf(v)
		t.Log(reflectType.String())
	}
}
