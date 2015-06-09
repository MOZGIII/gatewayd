package policy

import (
	"gatewayd/policy/sessionmanagement"
	"gatewayd/policy/tunnelaccess"
)

// Managers is a global policy manager.
var manager *ManagerType

// ManagerType is type for policy manager.
type ManagerType struct {
	SessionManagementPolicy sessionmanagement.Policy
	TunnelAccessPolicy      tunnelaccess.Policy
}

// SessionManagement policy returns current session management policy.
func SessionManagement() (sessionmanagement.Policy, error) {
	if manager == nil {
		return nil, ErrManagerNotReady
	}
	policy := manager.SessionManagementPolicy
	if policy == nil {
		return nil, ErrPolicyNotDefined
	}
	return policy, nil
}

// TunnelAccess policy returns current tunnel access policy.
func TunnelAccess() (tunnelaccess.Policy, error) {
	if manager == nil {
		return nil, ErrManagerNotReady
	}
	policy := manager.TunnelAccessPolicy
	if policy == nil {
		return nil, ErrPolicyNotDefined
	}
	return policy, nil
}
