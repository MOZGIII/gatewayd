package factorybase

import (
	"errors"
)

// ErrFailedCast can be used to report a failed cast.
var ErrFailedCast = errors.New("factorybase: failed to cast type in factory")
