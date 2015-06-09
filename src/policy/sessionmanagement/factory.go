package sessionmanagement

import (
	"gatewayd/pkg/factorybase"
)

// Factory stores data
var Factory = FactoryType{*factorybase.New()}

// FactoryProduceFunc is concrete produce func.
type FactoryProduceFunc func() Policy

// FactoryType is a concrete factory type.
type FactoryType struct {
	base factorybase.FactoryBase
}

// CreateByName is used to create a policy by it's name.
func (f *FactoryType) CreateByName(name string) (Policy, error) {
	val, err := f.base.CreateByName(name)
	if err != nil {
		return nil, err
	}
	policy, ok := val.(Policy)
	if !ok {
		return nil, factorybase.ErrFailedCast
	}
	return policy, nil
}

// Register allows to add new policies.
func (f *FactoryType) Register(name string, produceFunc FactoryProduceFunc) error {
	return f.base.Register(name, func() interface{} {
		return produceFunc()
	})
}

// Register is just an alias for Factory.Register
func Register(name string, produceFunc FactoryProduceFunc) error {
	return Factory.Register(name, produceFunc)
}
