package ajson

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/davecgh/go-spew/spew"
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
		unknowns map[string]any
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
			unknowns: map[string]any{
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
			unknowns: map[string]any{
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
			unknowns: map[string]any{
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
			unknowns: map[string]any{
				"sub.location": "USA",
				"location":     "USA",
			},
			expected: `{"name":"John","age":30,"sub":{"name":"Doe","location":"USA"},"location":"USA"}`,
		},
		"with unknown arrays": {
			sample: Sample{
				Name: "John",
				Age:  30,
			},
			unknowns: map[string]any{
				"location": []string{"USA", "UK"},
			},
			expected: `{"name":"John","age":30,"location":["USA","UK"]}`,
		},
		"with unknown maps": {
			sample: Sample{
				Name: "John",
				Age:  30,
			},
			unknowns: map[string]any{
				"location": map[string]string{
					"country": "USA",
					"city":    "New York",
				},
			},
			expected: `{"name":"John","age":30,"location":{"city":"New York","country":"USA"}}`,
		},
		"nil map[string]any": {
			sample: Sample{
				Name: "John",
				Age:  30,
			},
			unknowns: nil,
			expected: `{"name":"John","age":30}`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			data, err := Marshal(test.sample, test.unknowns)
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

func Test_recurseMap(t *testing.T) {
	tests := map[string]struct {
		m        map[string]any
		path     []string
		value    any
		expected map[string]any
	}{
		"simple": {
			m: map[string]any{
				"name": "John",
				"age":  30,
			},
			path:  []string{"location"},
			value: "USA",
			expected: map[string]any{
				"name":     "John",
				"age":      30,
				"location": "USA",
			},
		},
		"nested": {
			m: map[string]any{
				"name": "John",
				"age":  30,
			},
			path:  []string{"location", "country"},
			value: "USA",
			expected: map[string]any{
				"name": "John",
				"age":  30,
				"location": map[string]any{
					"country": "USA",
				},
			},
		},
		"nested with existing": {
			m: map[string]any{
				"name": "John",
				"age":  30,
				"location": map[string]any{
					"city": "New York",
				},
			},
			path:  []string{"location", "country"},
			value: "USA",
			expected: map[string]any{
				"name": "John",
				"age":  30,
				"location": map[string]any{
					"city":    "New York",
					"country": "USA",
				},
			},
		},
		"nested with existing and array": {
			m: map[string]any{
				"name": "John",
				"age":  30,
				"location": map[string]any{
					"city": "New York",
				},
			},
			path:  []string{"location", "country"},
			value: []string{"USA", "UK"},
			expected: map[string]any{
				"name": "John",
				"age":  30,
				"location": map[string]any{
					"city":    "New York",
					"country": []string{"USA", "UK"},
				},
			},
		},
		"nested with existing and map": {
			m: map[string]any{
				"name": "John",
				"age":  30,
				"location": map[string]any{
					"city": "New York",
				},
			},
			path: []string{"location", "country"},
			value: map[string]string{
				"country": "USA",
				"city":    "New York",
			},
			expected: map[string]any{
				"name": "John",
				"age":  30,
				"location": map[string]any{
					"city": "New York",
					"country": map[string]string{
						"country": "USA",
						"city":    "New York",
					},
				},
			},
		},
		"nested with existing and map and array": {
			m: map[string]any{
				"name": "John",
				"age":  30,
				"location": map[string]any{
					"city": "New York",
				},
			},
			path: []string{"location", "country"},
			value: map[string]any{
				"country": "USA",
				"city":    "New York",
				"places":  []string{"USA", "UK"},
			},
			expected: map[string]any{
				"name": "John",
				"age":  30,
				"location": map[string]any{
					"city": "New York",
					"country": map[string]any{
						"country": "USA",
						"city":    "New York",
						"places":  []string{"USA", "UK"},
					},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			recurseMap(test.m, test.path, test.value)

			delta := cmp.Diff(test.m, test.expected)
			if delta != "" {
				spew.Dump(test.m)
				spew.Dump(test.expected)
				t.Fatal(delta)
			}

			t.Logf("data: %s", test.m)
		})
	}
}

func TestUnmarshalJSON(t *testing.T) {
	tests := map[string]struct {
		data     string
		expected Sample
		unknowns map[string]any
	}{
		"simple": {
			data:     `{"name":"John","age":30}`,
			expected: Sample{Name: "John", Age: 30},
		},
		"simple with unknowns": {
			data:     `{"name":"John","age":30,"location":"USA"}`,
			expected: Sample{Name: "John", Age: 30},
			unknowns: map[string]any{"location": "USA"},
		},
		"with sub": {
			data:     `{"name":"John","age":30,"sub":{"name":"Doe"}}`,
			expected: Sample{Name: "John", Age: 30, Sub: &SubSample{Name: "Doe"}},
		},
		"with sub and unknowns": {
			data:     `{"name":"John","age":30,"sub":{"name":"Doe"},"location":"USA"}`,
			expected: Sample{Name: "John", Age: 30, Sub: &SubSample{Name: "Doe"}},
			unknowns: map[string]any{"location": "USA"},
		},
		"with sub-unknowns": {
			data:     `{"name":"John","age":30,"sub":{"name":"Doe","location":"USA"}}`,
			expected: Sample{Name: "John", Age: 30, Sub: &SubSample{Name: "Doe"}},
			unknowns: map[string]any{"sub.location": "USA"},
		},
		"with sub-unknowns and unknowns": {
			data:     `{"name":"John","age":30,"sub":{"name":"Doe","location":"USA"},"location":"USA"}`,
			expected: Sample{Name: "John", Age: 30, Sub: &SubSample{Name: "Doe"}},
			unknowns: map[string]any{"sub.location": "USA", "location": "USA"},
		},
		"with unknown arrays": {
			data:     `{"name":"John","age":30,"sub":{"name":"Doe","location":"USA"},"location":"USA","emails":["test@example.com","test2@example.com"]}`,
			expected: Sample{Name: "John", Age: 30, Sub: &SubSample{Name: "Doe"}},
			unknowns: map[string]any{
				"sub.location": "USA",
				"location":     "USA",
				"emails":       []string{"test@example.com", "test2@example.com"},
			},
		},
		"with unknown maps": {
			data:     `{"name":"John","age":30,"sub":{"name":"Doe","location":"USA"},"location":"USA","emails":{"test":"test@example.com"}}`,
			expected: Sample{Name: "John", Age: 30, Sub: &SubSample{Name: "Doe"}},
			unknowns: map[string]any{
				"sub.location": "USA",
				"location":     "USA",
				"emails":       map[string]string{"test": "test@example.com"},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sample, unknowns, err := Unmarshal[Sample]([]byte(test.data))
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

type Sample2 struct {
	Name string
	Age  int
	Sub  *SubSample `json:",omitempty"`
}

func Test_Unmarshal2(t *testing.T) {
	tests := map[string]struct {
		data     string
		expected Sample2
		unknowns map[string]any
	}{
		"simple": {
			data:     `{"name":"John","age":30}`,
			expected: Sample2{Name: "John", Age: 30},
		},
		"simple with unknowns": {
			data:     `{"name":"John","age":30,"location":"USA"}`,
			expected: Sample2{Name: "John", Age: 30},
			unknowns: map[string]any{"location": "USA"},
		},
		"with sub": {
			data:     `{"name":"John","age":30,"sub":{"name":"Doe"}}`,
			expected: Sample2{Name: "John", Age: 30, Sub: &SubSample{Name: "Doe"}},
		},
		"with sub and unknowns": {
			data:     `{"name":"John","age":30,"sub":{"name":"Doe"},"location":"USA"}`,
			expected: Sample2{Name: "John", Age: 30, Sub: &SubSample{Name: "Doe"}},
			unknowns: map[string]any{"location": "USA"},
		},
		"with sub-unknowns": {
			data:     `{"name":"John","age":30,"sub":{"name":"Doe","location":"USA"}}`,
			expected: Sample2{Name: "John", Age: 30, Sub: &SubSample{Name: "Doe"}},
			unknowns: map[string]any{"sub.location": "USA"},
		},
		"with sub-unknowns and unknowns": {
			data:     `{"name":"John","age":30,"sub":{"name":"Doe","location":"USA"},"location":"USA"}`,
			expected: Sample2{Name: "John", Age: 30, Sub: &SubSample{Name: "Doe"}},
			unknowns: map[string]any{"sub.location": "USA", "location": "USA"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			sample, unknowns, err := Unmarshal[Sample2]([]byte(test.data))
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

func Test_Marshal2(t *testing.T) {
	tests := map[string]struct {
		data     Sample2
		expected string
		unknowns map[string]any
	}{
		"simple": {
			data:     Sample2{Name: "John", Age: 30},
			expected: `{"Name":"John","Age":30}`,
		},
		"simple with unknowns": {
			data:     Sample2{Name: "John", Age: 30},
			expected: `{"Name":"John","Age":30,"location":"USA"}`,
			unknowns: map[string]any{"location": "USA"},
		},
		"with sub": {
			data:     Sample2{Name: "John", Age: 30, Sub: &SubSample{Name: "Doe"}},
			expected: `{"Name":"John","Age":30,"Sub":{"name":"Doe"}}`,
		},
		"with sub and unknowns": {
			data:     Sample2{Name: "John", Age: 30, Sub: &SubSample{Name: "Doe"}},
			expected: `{"Name":"John","Age":30,"Sub":{"name":"Doe"},"location":"USA"}`,
			unknowns: map[string]any{"location": "USA"},
		},
		"with sub-unknowns": {
			data:     Sample2{Name: "John", Age: 30, Sub: &SubSample{Name: "Doe"}},
			expected: `{"Name":"John","Age":30,"Sub":{"name":"Doe","location":"USA"}}`,
			unknowns: map[string]any{"Sub.location": "USA"},
		},
		"with sub-unknowns and unknowns": {
			data:     Sample2{Name: "John", Age: 30, Sub: &SubSample{Name: "Doe"}},
			expected: `{"Name":"John","Age":30,"Sub":{"name":"Doe","location":"USA"},"location":"USA"}`,
			unknowns: map[string]any{"Sub.location": "USA", "location": "USA"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			data, err := Marshal[Sample2](test.data, test.unknowns)
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

func BenchmarkMarshalJSON(b *testing.B) {
	sample := Sample{
		Name: "John",
		Age:  30,
		Sub: &SubSample{
			Name: "Doe",
		},
	}

	for i := 0; i < b.N; i++ {
		_, err := Marshal(sample, map[string]any{
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
		_, _, err := Unmarshal[Sample](data)
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

		unknowns := map[string]any{
			"location": stripCtlFromUTF8(location),
			"email":    stripCtlFromUTF8(email),
		}

		data, err := Marshal(s, unknowns)
		if err != nil {
			t.Fatalf("failed to marshal: %v", err)
		}

		decodedSample, decodedUnknowns, err := Unmarshal[Sample](data)
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
