package main

import (
	"testing"

	"github.com/melwyn95/code-generator/cmd/common"
)

func TestMarshallJSON(t *testing.T) {

	profile := common.Profile{
		Name:       "Melwyn Saldanha",
		Experience: 2,
		Hobbies:    []string{"Solve rubix cubes", "Watch movies"},
		Social: map[string]string{
			"github":  "https://github.com/melwyn95",
			"twitter": "https://twitter.com/MelwynSaldanha",
		},
	}

	want := `{"name":"Melwyn Saldanha","experience":2,"hobbies":["Solve rubix cubes","Watch movies"],"github":"https://github.com/melwyn95","twitter":"https://twitter.com/MelwynSaldanha"}`
	gotbytes, err := profile.MarshalJSON()
	got := string(gotbytes)

	if err != nil {
		t.Fatal("Error in reflection.MarshallJSON", err)
	}

	if got != want {
		t.Errorf("Expected: %s \nGot: %s", want, got)
	}
}
