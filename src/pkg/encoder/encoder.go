package encoder

import (
	"encoding/json"
)

// Encoder for abstract encoder.
type Encoder interface {
	Encode(v interface{}) ([]byte, error)
}

// Must panics on errors.
func Must(data []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return data
}

// JSONEncoder is a martini-compatible JSON encoder without filtering.
type JSONEncoder struct {
	PrettyPrint bool
}

// Encode does the trick.
func (e JSONEncoder) Encode(obj interface{}) ([]byte, error) {
	if e.PrettyPrint {
		return json.MarshalIndent(obj, "", "    ")
	}
	return json.Marshal(obj)
}
