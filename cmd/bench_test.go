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
	Social: map[string]string{
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
