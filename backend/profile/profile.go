package profile

// Profile holds data that control how sessions are executed.
type Profile struct {
	Name       string            `json:"name"`   // name uniquely identifies profile
	DriverName string            `json:"driver"` // driver name represents the driver used
	Params     map[string]string `json:"params"` // key/value pairs that are passed around for the system to use
}
