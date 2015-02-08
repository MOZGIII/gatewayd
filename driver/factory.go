package driver

import (
	"fmt"
)

// Factory makes drivers for sessions.

var (
	driverRegistry = make(map[string]FactoryProduceFunc)
)

// FactoryProduceFunc creates a driver of partilucar type.
type FactoryProduceFunc func() Driver

// CreateByName uses factory to create driver by it's name.
func CreateByName(driverName string) (Driver, error) {
	f, exists := driverRegistry[driverName]
	if !exists {
		return nil, fmt.Errorf("driver: no such driver %q found", driverName)
	}
	driver := f()
	return driver, nil
}

// FactoryRegister must be used to make factory to be able to produce
// a particular type. Used per-driver in init func.
func FactoryRegister(driverName string, driverCreationFunc FactoryProduceFunc) error {
	if _, exists := driverRegistry[driverName]; exists {
		return fmt.Errorf("driver: driver with name %q is already registered", driverName)
	}

	driverRegistry[driverName] = driverCreationFunc
	return nil
}

// Registry return the factroy driver registry, used for debugging purposes.
func Registry() map[string]FactoryProduceFunc {
	return driverRegistry
}
