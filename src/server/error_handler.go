package server

import (
	"log"
	"net/http"

	"gatewayd/utils"
)

// ErrorWrapper wraps net.http callbacks with error handling
func ErrorWrapper(f func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			switch err := err.(type) {
			case utils.MethodNotAllowedError:
				http.Error(w, "Method not allowed", 405)
			case utils.HTTPError:
				http.Error(w, err.Error(), err.Code())
				log.Printf("handling %q: %v", r.RequestURI, err)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
				log.Printf("handling %q: %v", r.RequestURI, err)
			}
		}
	}
}
