package profile

import (
	"fmt"
)

// Profile holds data that control how sessions are executed.
type Profile struct {
	Name       string                 `json:"name"`   // name uniquely identifies profile
	DriverName string                 `json:"driver"` // driver name represents the driver used
	Params     map[string]interface{} `json:"params"` // key/value pairs that are passed around for the system to use
}

// GetParam return param for specified key.
func (p *Profile) GetParam(key string) (interface{}, error) {
	val, ok := p.Params[key]
	if !ok {
		return nil, fmt.Errorf("profile: param not found for key %q", key) // FIXME
	}

	return val, nil
}

// GetParams return param
func (p *Profile) GetParams() map[string]interface{} {
	return p.Params
}
