package common

//go:generate json-generator

type Profile struct {
	Name       string            `json:"name"`
	Experience int               `json:"experience"`
	Hobbies    []string          `json:"hobbies"`
	Social     map[string]string `json:"social"`
}
