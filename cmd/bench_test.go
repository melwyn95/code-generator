package main

import (
	"testing"

	"github.com/melwyn95/code-generator/cmd/common"
	reflection "github.com/melwyn95/code-generator/cmd/json-reflection"
)

var profile common.Profile = common.Profile{
	Name:       "Melwyn Saldanha",
	Experience: 2,
	Hobbies:    []string{"Solve rubix cubes", "Watch movies"},
	RandomStuff: map[string]string{
		"github":  "https://github.com/melwyn95",
		"twitter": "https://twitter.com/MelwynSaldanha",
	},
}

func BenchmarkReflectionJSONMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reflection.MarshallJSON(profile)
	}
}

func BenchmarkGeneratedJSONMarshall(b *testing.B) {
	for i := 0; i < b.N; i++ {
		profile.MarshalJSON()
	}
}

func TestReflectionAndGeneratedCode(t *testing.T) {
	reflectionJSON, _ := reflection.MarshallJSON(profile)
	generatorJSON, _ := profile.MarshalJSON()

	if string(reflectionJSON) != string(generatorJSON) {
		t.Errorf("Reflection: %s \n Generator: %s", string(reflectionJSON), string(generatorJSON))
	}
}
