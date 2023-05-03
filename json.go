package ajson

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-json-experiment/json"
	"github.com/trivago/tgo/tcontainer"
	"go.devnw.com/structs"
)

type MMap tcontainer.MarshalMap

func MarshalJSON[T comparable](t T, mm MMap) ([]byte, error) {
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
	for key, v := range m {
		if key == path[0] {
			subMap, ok := v.(map[string]any)
			if !ok {
				continue
			}

			recurseMap(subMap, path[1:], value)
		}
	}
}

func UnmarshalJSON[T comparable](data []byte) (T, MMap, error) {
	mm := tcontainer.NewMarshalMap()

	var t T
	err := json.Unmarshal(data, &t)
	if err != nil {
		return t, MMap(mm), err
	}

	err = json.Unmarshal(data, &mm)
	if err != nil {
		return t, MMap(mm), err
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
			return t, MMap(mm),
				fmt.Errorf("no tags options found for %s", field.Name)
		}

		// the first one is the json tag
		key := options[0]
		if _, okay := mm.Value(key); okay {
			delete(mm, key)
		}
	}

	return t, MMap(mm), nil
}
