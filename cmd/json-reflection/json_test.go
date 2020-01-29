package reflection

import "testing"

func TestMarshallJSON(t *testing.T) {
	type Profile struct {
		Name        string            `json:"name"`
		Experience  float64           `json:"experience"`
		Hobbies     []string          `json:"hobbies"`
		RandomStuff map[string]string `json:"random_stuff"`
	}
	profile := Profile{
		Name:       "Melwyn Saldanha",
		Experience: 2.6,
		Hobbies:    []string{"Solve rubix cubes", "Watch movies"},
		RandomStuff: map[string]string{
			"github":  "https://github.com/melwyn95",
			"twitter": "https://twitter.com/MelwynSaldanha",
		},
	}

	want := `{"name":"Melwyn Saldanha","experience":"2.6","hobbies":["Solve rubix cubes","Watch movies"],"github":"https://github.com/melwyn95","twitter":"https://twitter.com/MelwynSaldanha"}`
	gotbytes, err := MarshallJSON(profile)
	got := string(gotbytes)

	if err != nil {
		t.Fatal("Error in reflection.MarshallJSON", err)
	}

	if got != want {
		t.Errorf("Expected: %s \nGot: %s", want, got)
	}
}
