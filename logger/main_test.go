package main_test

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

type Hoge struct {
	Foo  string `json:"foo"`
	Bar  string `json:"bar"`
	Baz  string `json:"baz"`
	Fuga Fuga   `json:"fuga"`
}

type Fuga struct {
	Piyo string `json:"piyo"`
}

var (
	hoge   = Hoge{Foo: "foo", Bar: "bar", Baz: "baz", Fuga: Fuga{Piyo: "piyo"}}
	expect = map[string]interface{}{"foo": "foo", "bar": "bar", "baz": "baz", "fuga": map[string]interface{}{"piyo": "piyo"}}
)

// old
func MarshalUnmarshal(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	b, _ := json.Marshal(data)
	json.Unmarshal(b, &result)

	return result
}

// new
func Reflect(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	typ := reflect.TypeOf(data)
	// should check data is struct or not
	val := reflect.ValueOf(data)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		vi := val.FieldByName(field.Name).Interface()
		// if field is struct, convert recursively
		if field.Type.Kind() == reflect.Struct {
			vi = Reflect(vi)
		}
		if tag, ok := field.Tag.Lookup("json"); ok {
			result[tag] = vi
			continue
		}
		result[strings.ToLower(field.Name)] = vi
	}
	return result
}

/* test old and new function is same*/

func TestMarshalUnmarshal(t *testing.T) {
	out := MarshalUnmarshal(hoge)
	if !reflect.DeepEqual(out, expect) {
		t.Errorf("%s != %s", out, expect)
	}
}

func TestReflect(t *testing.T) {
	out := Reflect(hoge)
	for k, v := range out {
		t.Logf("%s:%s : %s:%s", k, reflect.TypeOf(k), v, reflect.TypeOf(v))
	}
	if !reflect.DeepEqual(out, expect) {
		t.Errorf("%s != %s", out, expect)
	}
}

/* benchmark */

func BenchmarkMarshalUnmarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MarshalUnmarshal(hoge)
	}
}

func BenchmarkReflect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Reflect(hoge)
	}
}
