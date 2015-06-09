package policy

import (
	"errors"

	"gatewayd/config"

	"gatewayd/policy/sessionmanagement"
	"gatewayd/policy/tunnelaccess"
)

// LoadFromConfig initializes global policy manager
// from the config params passed.
// This can be called only once.
func LoadFromConfig(c *config.PolicyConfig) error {
	if manager != nil {
		return errors.New("policy: policy manager is already initialized")
	}

	// Session management policy
	sessionManagementPolicy, err := sessionmanagement.Factory.CreateByName(c.SessionManagementPolicyName)
	if err != nil {
		return err
	}

	// Tunnel access policy
	tunnelAccessPolicy, err := tunnelaccess.Factory.CreateByName(c.TunnelAccessPolicyName)
	if err != nil {
		return err
	}

	manager = &ManagerType{
		sessionManagementPolicy,
		tunnelAccessPolicy,
	}
	return nil
}
