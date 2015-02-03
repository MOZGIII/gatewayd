package api

import (
	"encoding/json"
	"log"
	"net/http"

	"gatewayd/backend"
	"gatewayd/utils"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/encoder"
)

// TODO: optimize code in here, we need some better apprach in general
// maybe...

// SessionByKey responds to the API request to the session by key
func SessionByKey(params martini.Params, enc encoder.Encoder) (int, []byte) {
	key, ok := params["key"]
	if !ok {
		return http.StatusBadRequest, []byte("No key passed")
	}
	log.Printf("Requesting session info for key %q", key)
	backend.Control.DebugChan() <- utils.NewEmpty()

	ansc := make(chan *backend.Session)
	backend.Control.SyncChannel() <- func(sm *backend.SessionsManager) {
		ansc <- sm.SessionByKey(key)
	}
	val := <-ansc
	// Seems like accessing session in here is not safe cause it might be used
	// from other goroutines after sync has passed.

	if val == nil {
		log.Println("Requested session not found with key", key)
		return http.StatusNotFound, []byte("No such session")
	}

	return http.StatusOK, encoder.Must(enc.Encode(val.Export()))
}

// SessionByToken responds to the API request to the session by token
func SessionByToken(params martini.Params, enc encoder.Encoder) (int, []byte) {
	token, ok := params["token"]
	if !ok {
		return http.StatusBadRequest, []byte("No token passed")
	}
	log.Printf("Requesting session info for token %q", token)
	backend.Control.DebugChan() <- utils.NewEmpty()

	ansc := make(chan *backend.Session)
	backend.Control.SyncChannel() <- func(sm *backend.SessionsManager) {
		ansc <- sm.SessionByToken(token)
	}
	val := <-ansc

	if val == nil {
		log.Println("Requested session not found with token", token)
		return http.StatusNotFound, []byte("No such session")
	}

	return http.StatusOK, encoder.Must(enc.Encode(val.Export()))
}

// CreateSession creates a new session
func CreateSession(params martini.Params, enc encoder.Encoder, req *http.Request) (int, []byte) {
	decoder := json.NewDecoder(req.Body)
	var t struct {
		Key string `json:"key"`
	}
	if err := decoder.Decode(&t); err != nil {
		log.Println(err)
		return http.StatusBadRequest, []byte("Wrong request")
	}
	log.Println("Creating session with key", t.Key)

	type ansv struct {
		session *backend.Session
		err     error
	}
	ansc := make(chan ansv)
	backend.Control.SyncChannel() <- func(sm *backend.SessionsManager) {
		session, err := sm.CreateSession(t.Key)
		ansc <- ansv{session, err}
	}
	val := <-ansc
	log.Println(val.session)

	if val.err != nil {
		log.Println(val.err)
		return http.StatusInternalServerError, []byte("Unable to create session")
	}

	return http.StatusCreated, encoder.Must(enc.Encode(val.session.Export()))
}

// GetOrCreate
// func GetOrCreateSession(rw http.ResponseWriter, req *http.Request) error {
// 	return nil
// }
