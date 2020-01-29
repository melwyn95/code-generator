package common

import (
	"bytes"
	"fmt"
)

func (p *Profile) MarshalJSON() {
	var json bytes.Buffer
	json.WriteString("{")
	json.WriteString(`"name":`)
	json.WriteString(fmt.Sprintf("", p.Name) + `,`)
	json.WriteString(`"experience":`)
	json.WriteString(fmt.Sprintf("", p.Experience) + `,`)
	json.WriteString(`[`)
	for i := range p.Hobbies {
		json.WriteString(fmt.Sprintf("", p.Hobbies[i]) + `,`)
	}
	json.WriteString(`],`)
	for k, v := range p.RandomStuff {
		json.WriteString(`"` + k + `"` + `:` + `"` + v + `",`)
	}
	json.WriteString("}")
}
