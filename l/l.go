package l

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

type Logger struct{}

var Log Logger = Logger{}

func (i Logger) Log(input interface{}) {

}

var LogFunc = log.Println
var FatalFunc = log.Fatal

func F(i interface{}) {
	FatalFunc(Parse(i))
}
func Parse(i interface{}) interface{} {
	typeOf := reflect.TypeOf(i)
	if typeOf.Kind() == reflect.Map {
		b, err := json.MarshalIndent(i, "", "  ")
		if err != nil {
			return i
		}
		return string(b)
	}
	if typeOf.Kind() == reflect.Struct {
		b, err := json.MarshalIndent(i, "", "  ")
		if err != nil {
			return i
		}
		name := typeOf.Name()
		result := fmt.Sprintf("<%v>%v", name, string(b))
		return result
	}
	if typeOf.Kind() == reflect.Slice {
		v := reflect.ValueOf(i)
		result := "["
		for i := 0; i < v.Len(); i++ {
			val := v.Index(i)
			result += fmt.Sprintf("\n%v,", Parse(val.Interface()))
		}
		result = result[:len(result)-1] + "\n"
		result += "]"
		return result
	}
	if typeOf.Kind() == reflect.Array {
		v := reflect.ValueOf(i)
		result := "["
		for i := 0; i < v.Len(); i++ {
			val := v.Index(i)
			result += fmt.Sprintf("%v,", Parse(val.Interface()))
		}
		result += "]"
		return result
	}
	return i
}
func L(i interface{}) {
	LogFunc(Parse(i))
}
