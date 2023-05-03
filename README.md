# ajson
--
    import "."

Package ajson provides a way to marshal and unmarshal JSON with unknown fields.

[![Build & Test Action
Status](https://github.com/devnw/ajson/actions/workflows/build.yml/badge.svg)](https://github.com/devnw/ajson/actions)
[![Go Report
Card](https://goreportcard.com/badge/go.devnw.com/ajson)](https://goreportcard.com/report/go.devnw.com/ajson)
[![codecov](https://codecov.io/gh/devnw/ajson/branch/main/graph/badge.svg)](https://codecov.io/gh/devnw/ajson)
[![Go
Reference](https://pkg.go.dev/badge/go.devnw.com/ajson.svg)](https://pkg.go.dev/go.devnw.com/ajson)
[![License: Apache
2.0](https://img.shields.io/badge/license-Apache-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![PRs
Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](http://makeapullrequest.com)

## Installation

```bash

go get -u go.devnw.com/ajson

```

## Usage

#### func  Marshal

```go
func Marshal[T comparable](t T, mm MMap) ([]byte, error)
```
MarshalJSON marshals the given struct to json and then merges the unknown fields
into the json from the MMap object

Example usage:

    type Sample struct {
    	Name string     `json:"name"`
    	Age  int        `json:"age"`
    	Sub  *SubSample `json:"sub,omitempty"`
    }

    type SubSample struct {
    	Name string `json:"name"`
    }

    func main() {
    	sample := Sample{
    		Name: "John",
    		Age:  30,
    	}

    	unknowns := MMap{
    		"location": "USA",
    	}

    	data, err := MarshalJSON(sample, unknowns)
    	if err != nil {
    		panic(err)
    	}

    	fmt.Println(string(data))
    }

    // Output:
    // {"name":"John","age":30,"location":"USA"}

Example with embeded unknown and custom marshaler:

    type Sample struct {
    	Name 		string
    	Age  		int
    	Unknowns	MMap
    }

    func (s Sample) MarshalJSON() ([]byte, error) {
    	return MarshalJSON(struct {
    		ID   string `json:"id"`
    		Name string `json:"name"`
    	}{
    		ID:   t.ID,
    		Name: t.Name,
    	}, t.Unknowns)
    }

#### type MMap

```go
type MMap tcontainer.MarshalMap
```


#### func  Unmarshal

```go
func Unmarshal[T comparable](data []byte) (T, MMap, error)
```
UnmarshalJSON unmarshals the given json into the given struct and then returns
the unknown fields as a MMap object.
