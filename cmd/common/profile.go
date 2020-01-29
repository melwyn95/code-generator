package common

//go:generate json-generator

type Profile struct {
	Name        string            `json:"name"`
	Experience  float64           `json:"experience"`
	Hobbies     []string          `json:"hobbies"`
	RandomStuff map[string]string `json:"random_stuff"`
}
