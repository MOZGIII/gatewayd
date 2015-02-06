package api

import (
	"encoding/json"
	"log"
	"net/http"

	"gatewayd/backend/control"
	"gatewayd/backend/session"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/encoder"
)

// API talks to session manager to get sessions.
// Remember sessions live each in it's own goroutine.

// SessionByToken responds to the API request to the session by token
func SessionByToken(params martini.Params, enc encoder.Encoder) (int, []byte) {
	token, ok := params["token"]
	if !ok {
		return http.StatusBadRequest, []byte("No token passed")
	}
	log.Printf("Requesting session info for token %q", token)

	_, err := control.FixmeContextExport().SessionManager.SessionByToken(token)
	if err != nil {
		log.Println("Requested session not found with token", token)
		log.Println(err)
		return http.StatusNotFound, []byte("No such session")
	}

	info := struct {
		Status string `json:"tmp_status"`
	}{"online"}
	return http.StatusOK, encoder.Must(enc.Encode(info))
}

// CreateSession creates a new session and returns it's data.
// With this API call, you should be able to ...
// When this call returns, session is just set to be spawning and may not
// be done with spawning yet (so it may be unconnectable).
func CreateSession(params martini.Params, enc encoder.Encoder, req *http.Request) (int, []byte) {
	decoder := json.NewDecoder(req.Body)
	var t struct {
		ProfileName   string            `json:"profile"`
		DynamicParams map[string]string `json:"params"` // not those params that are in profile! these are different!
	}
	if err := decoder.Decode(&t); err != nil {
		log.Println(err)
		return http.StatusBadRequest, []byte("Wrong request")
	}
	log.Println("Creating session for request", t)

	session := session.New(nil, nil)

	token, err := control.FixmeContextExport().SessionManager.Register(session)
	if err != nil {
		session.Terminate()
		log.Println(err)
		return http.StatusInternalServerError, []byte("Unable to create session")
	}
	log.Printf("Session created and registered (token %q)", token)

	result := struct {
		Token string `json:"token"`
	}{token}
	return http.StatusCreated, encoder.Must(enc.Encode(result))
}
