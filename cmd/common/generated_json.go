package common

import (
	"bytes"
	"strconv"
)

//DO NOT EDIT. This code is auto-generated using go:generate json-generator

func truncateLastComma(buff *bytes.Buffer) {
	if buff != nil {
		bufflen := buff.Len()
		if bufflen > 0 && string(buff.Bytes()[bufflen-1:]) == "," {
			buff.Truncate(bufflen - 1)
		}
	}
}
func (p *Profile) MarshalJSON() ([]byte, error) {
	var json bytes.Buffer
	json.WriteString("{")
	json.WriteString(`"name":`)
	json.WriteString(`"` + p.Name + `"`)
	json.WriteString(",")
	json.WriteString(`"experience":`)
	json.WriteString(strconv.Itoa(p.Experience))
	json.WriteString(",")
	json.WriteString(`"hobbies":`)
	json.WriteString(`[`)
	for i := range p.Hobbies {
		json.WriteString(`"` + p.Hobbies[i] + `"`)
		json.WriteString(",")
	}
	truncateLastComma(&json)
	json.WriteString(`],`)
	for k, v := range p.RandomStuff {
		json.WriteString(`"` + k + `"` + `:` + `"` + v + `",`)
	}
	truncateLastComma(&json)
	json.WriteString("}")
	truncateLastComma(&json)
	return json.Bytes(), nil
}
