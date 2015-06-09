package policy

import (
	"errors"
)

// ErrManagerNotReady happens when policy manager is
// accessed but was not initialized yet.
var ErrManagerNotReady = errors.New("policy: policy manager is not ready")

// ErrPolicyNotDefined is for when the policy that is
// not configured is accessed.
var ErrPolicyNotDefined = errors.New("policy: policy is not defined")
