package common

//go:generate json-generator

type Profile struct {
	Name        string            `json:"name"`
	Experience  int               `json:"experience"`
	Hobbies     []string          `json:"hobbies"`
	RandomStuff map[string]string `json:"random_stuff"`
}
