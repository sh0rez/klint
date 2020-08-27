package klint

import (
	"fmt"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

// Get returns any nested field of the manifest. It is invoked as m.Get("metadata", "name")
func GetField(m manifest.Manifest, fieldspec ...interface{}) (*Got, error) {
	var last interface{} = map[string]interface{}(m)

	for _, f := range fieldspec {
		switch t := f.(type) {
		case string:
			msi, ok := last.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("cannot select map-key '%s' of '%T'. Expected 'map[string]interface{}'", f, last)
			}

			v, ok := msi[t]
			if !ok {
				return nil, fmt.Errorf("field '%s' does not exist", t)
			}

			last = v
		case int:
			slice, ok := last.([]interface{})
			if !ok {
				return nil, fmt.Errorf("type '%T' is not indexable", last)
			}

			if l := len(slice); t >= l {
				return nil, fmt.Errorf("Index %v exceeds slice length %v", t, l)
			}

			last = slice[t]
		default:
			return nil, fmt.Errorf("fieldspec must be only string (map keys) and int (slice indexes). Found %T", t)
		}
	}

	return &Got{last}, nil
}

// Got is the result of a manifest.Get call
type Got struct {
	data interface{}
}

// Interface returns the raw data
func (g Got) Interface() interface{} {
	return interface{}(g.data)
}

// Map casts the underlying data to a map[string]interface{}. Panics if of wrong type
func (g Got) Map() map[string]interface{} {
	return g.data.(map[string]interface{})
}

// Slice casts the underlying data to a []interface{}. Panics if of wrong type
func (g *Got) Slice() []interface{} {
	return g.data.([]interface{})
}
