package reflection

import (
	"bytes"
	"fmt"
	"reflect"
)

// START - Dummy tag inorder to ignore the tag while Marshaling
const START = "@@"

func truncateLastComma(buff *bytes.Buffer) {
	if buff != nil {
		bufflen := buff.Len()
		if bufflen > 0 && string(buff.Bytes()[bufflen-1:]) == "," {
			buff.Truncate(bufflen - 1)
		}
	}
}

func marshall(source interface{}, buff *bytes.Buffer, tag string) {
	value := reflect.ValueOf(source)
	t := reflect.TypeOf(source)

	switch t.Kind() {
	case reflect.Struct:
		buff.WriteString("{")
		for i := 0; i < t.NumField(); i++ {
			marshall(value.Field(i).Interface(), buff, t.Field(i).Tag.Get("json"))
		}
		truncateLastComma(buff)
		buff.WriteString("}")
	case reflect.Slice:
		buff.WriteString(fmt.Sprintf(`"%s":[`, tag))
		for i := 0; i < value.Len(); i++ {
			marshall(value.Index(i).Interface(), buff, START)
		}
		truncateLastComma(buff)
		buff.WriteString("],")
	case reflect.Map:
		keys := value.MapKeys()
		for _, key := range keys {
			val := value.MapIndex(key)
			buff.WriteString(fmt.Sprintf(`"%s":"%s",`, key, val))
		}
	case reflect.Float64, reflect.String:
		if tag == START {
			buff.WriteString(fmt.Sprintf(`"%s",`, value))
		} else {
			buff.WriteString(fmt.Sprintf(`"%s":"%v",`, tag, value))
		}
	}
}

func MarshallJSON(source interface{}) ([]byte, error) {
	var buff bytes.Buffer
	marshall(source, &buff, START)
	truncateLastComma(&buff)

	return buff.Bytes(), nil
}
