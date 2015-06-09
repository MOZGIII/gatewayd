package factorybase

import (
	"fmt"
)

// ProduceFunc creates an object of some kind.
type ProduceFunc func() interface{}

// FactoryBase is used to store factory
type FactoryBase struct {
	m map[string]ProduceFunc
}

// New creates and initialzies new factory.
func New() *FactoryBase {
	return &FactoryBase{
		make(map[string]ProduceFunc),
	}
}

// CreateByName uses factory to produce a new object by it's name.
// This can be used for internal operations.
func (f *FactoryBase) CreateByName(name string) (interface{}, error) {
	pf, exists := f.m[name]
	if !exists {
		return nil, fmt.Errorf("factorybase: no such procuder %q found", name)
	}
	return pf(), nil
}

// Register must be used to add producers to factory.
func (f *FactoryBase) Register(name string, produceFunc ProduceFunc) error {
	if _, exists := f.m[name]; exists {
		return fmt.Errorf("factorybase: producer with name %q is already registered", name)
	}

	f.m[name] = produceFunc
	return nil
}
