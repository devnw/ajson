package ajson

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"go.devnw.com/structs"
)

// MarshalJSON marshals the given struct to json and then merges
// the unknown fields into the json from the map[string]any object
//
// Example usage:
//
//	type Sample struct {
//		Name string     `json:"name"`
//		Age  int        `json:"age"`
//		Sub  *SubSample `json:"sub,omitempty"`
//	}
//
//	type SubSample struct {
//		Name string `json:"name"`
//	}
//
//	func main() {
//		sample := Sample{
//			Name: "John",
//			Age:  30,
//		}
//
//		unknowns := map[string]any{
//			"location": "USA",
//		}
//
//		data, err := MarshalJSON(sample, unknowns)
//		if err != nil {
//			panic(err)
//		}
//
//		fmt.Println(string(data))
//	}
//
//	// Output:
//	// {"name":"John","age":30,"location":"USA"}
//
// Example with embeded unknown and custom marshaler:
//
//	type Sample struct {
//		Name 		string
//		Age  		int
//		Unknowns	map[string]any
//	}
//
//	func (s Sample) MarshalJSON() ([]byte, error) {
//		return MarshalJSON(struct {
//			ID   string `json:"id"`
//			Name string `json:"name"`
//		}{
//			ID:   t.ID,
//			Name: t.Name,
//		}, t.Unknowns)
//	}
func Marshal[T comparable](t T, mm map[string]any) ([]byte, error) {
	s := structs.New(t)
	s.TagName = "json"

	m := s.Map()
	for key, value := range mm {
		recurseMap(m, strings.Split(key, "."), value)
	}

	return json.Marshal(m)
}

func recurseMap(m map[string]any, path []string, value any) {
	if len(path) == 1 {
		m[path[0]] = value
		return
	}

	// if there are more than one path, we need to find the
	// correct map to set the value
	// e.g. "sub.name" -> m["sub"].(map[string]any)["name"]
	found := false
	for key, v := range m {
		if key == path[0] {
			subMap, ok := v.(map[string]any)
			if !ok {
				continue
			}

			recurseMap(subMap, path[1:], value)
			found = true
		}
	}

	if !found {
		newSubMap := make(map[string]any)
		recurseMap(newSubMap, path[1:], value)
		m[path[0]] = newSubMap
	}
}

// UnmarshalJSON unmarshals the given json into the given struct
// and then returns the unknown fields as a map[string]any object.
//
// Example usage:
//
//	type Sample struct {
//		Name string     `json:"name"`
//		Age  int        `json:"age"`
//		Sub  *SubSample `json:"sub,omitempty"`
//	}
//
//	type SubSample struct {
//		Name string `json:"name"`
//	}
//
//	func main() {
//		data := []byte(`{"name":"John","age":30,"location":"USA"}`)
//
//		var sample Sample
//		unknowns, err := UnmarshalJSON(data, &sample)
//		if err != nil {
//			panic(err)
//		}
//
//		fmt.Println(sample)
//		fmt.Println(unknowns)
//	}
//
//	// Output:
//	// {John 30 &{ }}
//	// map[location:USA]
//
// Example with embeded unknown and custom unmarshaler:
//
//	type Sample struct {
//		Name 		string
//		Age  		int
//		Unknowns	map[string]any
//	}
//
//	func (s *Sample) UnmarshalJSON(data []byte) error {
//		var t struct {
//			ID   string `json:"id"`
//			Name string `json:"name"`
//		}
//
//		unknowns, err := UnmarshalJSON(data, &t)
//		if err != nil {
//			return err
//		}
//
//		s.Name = t.Name
//		s.Age = t.Age
//		s.Unknowns = unknowns
//
//		return nil
//	}
func Unmarshal[T comparable](data []byte) (T, map[string]any, error) {
	mm := map[string]any{}

	var t T
	err := json.Unmarshal(data, &t)
	if err != nil {
		return t, mm, err
	}

	err = json.Unmarshal(data, &mm)
	if err != nil {
		return t, mm, err
	}

	tpe := reflect.TypeOf(t)
	for i := 0; i < tpe.NumField(); i++ {
		field := tpe.Field(i)
		tagDetail := field.Tag.Get("json")
		if tagDetail == "" {
			// ignore if there are no tags
			continue
		}
		options := strings.Split(tagDetail, ",")

		if len(options) == 0 {
			return t, mm,
				fmt.Errorf("no tags options found for %s", field.Name)
		}

		// the first one is the json tag
		delete(mm, options[0])
	}

	return t, mm, nil
}
