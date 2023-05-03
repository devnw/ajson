package ajson

import (
	"strings"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/google/go-cmp/cmp"
)

type Sample struct {
	Name string     `json:"name"`
	Age  int        `json:"age"`
	Sub  *SubSample `json:"sub,omitempty"`
}

type SubSample struct {
	Name string `json:"name"`
}

func TestMarshalJSON(t *testing.T) {
	tests := map[string]struct {
		sample   Sample
		unknowns MMap
		expected string
	}{
		"simple": {
			sample: Sample{
				Name: "John",
				Age:  30,
			},
			expected: `{"name":"John","age":30}`,
		},
		"with sub": {
			sample: Sample{
				Name: "John",
				Age:  30,
				Sub: &SubSample{
					Name: "Doe",
				},
			},
			expected: `{"name":"John","age":30,"sub":{"name":"Doe"}}`,
		},
		"with unknowns": {
			sample: Sample{
				Name: "John",
				Age:  30,
			},
			unknowns: MMap{
				"location": "USA",
				"email":    "test@email.com",
			},
			expected: `{"name":"John","age":30,"location":"USA","email":"test@email.com"}`,
		},
		"with sub and unknowns": {
			sample: Sample{
				Name: "John",
				Age:  30,
				Sub: &SubSample{
					Name: "Doe",
				},
			},
			unknowns: MMap{
				"location": "USA",
			},
			expected: `{"name":"John","age":30,"sub":{"name":"Doe"},"location":"USA"}`,
		},
		"with sub-unknowns": {
			sample: Sample{
				Name: "John",
				Age:  30,
				Sub: &SubSample{
					Name: "Doe",
				},
			},
			unknowns: MMap{
				"sub.location": "USA",
			},
			expected: `{"name":"John","age":30,"sub":{"name":"Doe","location":"USA"}}`,
		},
		"with sub-unknowns and unknowns": {
			sample: Sample{
				Name: "John",
				Age:  30,
				Sub: &SubSample{
					Name: "Doe",
				},
			},
			unknowns: MMap{
				"sub.location": "USA",
				"location":     "USA",
			},
			expected: `{"name":"John","age":30,"sub":{"name":"Doe","location":"USA"},"location":"USA"}`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			data, err := MarshalJSON(test.sample, test.unknowns)
			if err != nil {
				t.Fatal(err)
			}

			delta, err := diff(string(data), test.expected)
			if err != nil {
				t.Fatal(err)
			}

			if delta != "" {
				t.Fatal(delta)
			}

			t.Logf("data: %s", data)
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := map[string]struct {
		data     string
		expected Sample
		unknowns MMap
	}{
		"simple": {
			data:     `{"name":"John","age":30}`,
			expected: Sample{Name: "John", Age: 30},
		},
		"simple with unknowns": {
			data:     `{"name":"John","age":30,"location":"USA"}`,
			expected: Sample{Name: "John", Age: 30},
			unknowns: MMap{"location": "USA"},
		},
		"with sub": {
			data:     `{"name":"John","age":30,"sub":{"name":"Doe"}}`,
			expected: Sample{Name: "John", Age: 30, Sub: &SubSample{Name: "Doe"}},
		},
		"with sub and unknowns": {
			data:     `{"name":"John","age":30,"sub":{"name":"Doe"},"location":"USA"}`,
			expected: Sample{Name: "John", Age: 30, Sub: &SubSample{Name: "Doe"}},
			unknowns: MMap{"location": "USA"},
		},
		"with sub-unknowns": {
			data:     `{"name":"John","age":30,"sub":{"name":"Doe","location":"USA"}}`,
			expected: Sample{Name: "John", Age: 30, Sub: &SubSample{Name: "Doe"}},
			unknowns: MMap{"sub.location": "USA"},
		},
		"with sub-unknowns and unknowns": {
			data:     `{"name":"John","age":30,"sub":{"name":"Doe","location":"USA"},"location":"USA"}`,
			expected: Sample{Name: "John", Age: 30, Sub: &SubSample{Name: "Doe"}},
			unknowns: MMap{"sub.location": "USA", "location": "USA"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sample, unknowns, err := UnmarshalJSON[Sample]([]byte(test.data))
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(sample, test.expected); diff != "" {
				t.Fatal(diff)
			}

			t.Logf("sample: %+v", sample)
			t.Logf("unknowns: %+v", unknowns)
		})
	}
}

func BenchmarkMarshalJSON(b *testing.B) {
	sample := Sample{
		Name: "John",
		Age:  30,
		Sub: &SubSample{
			Name: "Doe",
		},
	}

	for i := 0; i < b.N; i++ {
		_, err := MarshalJSON(sample, MMap{
			"location": "USA",
			"email":    "test@email.com",
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshalJSON(b *testing.B) {
	data := []byte(`{"name":"John","age":30,"sub":{"name":"Doe"},"location":"USA"}`)

	for i := 0; i < b.N; i++ {
		_, _, err := UnmarshalJSON[Sample](data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func FuzzMarshalJSON(f *testing.F) {
	f.Fuzz(func(t *testing.T, name string, age int, location, email string) {
		s := Sample{
			Name: stripCtlFromUTF8(name),
			Age:  age,
		}

		unknowns := MMap{
			"location": stripCtlFromUTF8(location),
			"email":    stripCtlFromUTF8(email),
		}

		data, err := MarshalJSON(s, unknowns)
		if err != nil {
			t.Fatalf("failed to marshal: %v", err)
		}

		decodedSample, decodedUnknowns, err := UnmarshalJSON[Sample](data)
		if err != nil {
			t.Fatalf("failed to unmarshal: %v", err)
		}

		if decodedSample.Name != s.Name {
			t.Fatalf("name mismatch: %s != %s", decodedSample.Name, s.Name)
		}

		if decodedSample.Age != s.Age {
			t.Fatalf("age mismatch: %d != %d", decodedSample.Age, s.Age)
		}

		if len(decodedUnknowns) != len(unknowns) {
			t.Fatalf("unknowns mismatch: %d != %d", len(decodedUnknowns), len(unknowns))
		}

		for k, v := range unknowns {
			if decodedUnknowns[k] != v {
				t.Fatalf("unknowns mismatch: %s != %s", decodedUnknowns[k], v)
			}
		}
	})
}

func diff(s1, s2 string) (string, error) {
	var o1, o2 any

	err := json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return "", err
	}

	return cmp.Diff(o1, o2), nil
}

func stripCtlFromUTF8(str string) string {
	return strings.Map(func(r rune) rune {
		if r >= 32 && r != 127 {
			return r
		}
		return -1
	}, str)
}
