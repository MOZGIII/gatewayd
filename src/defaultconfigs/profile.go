package defaultconfigs

import (
	"io"
	"strings"
)

// GetProfileJSON provides the sample profile file JSON.
func GetProfileJSON() io.Reader {
	return strings.NewReader(`{
		"profiles": [
			{
				"name": "test",
				"driver": "localexec",
				"params": {
					"command": {
						"name": "gatewayd-session-test",
						"args": ["test"]
					}
				}
			}
		]
	}`)
}
