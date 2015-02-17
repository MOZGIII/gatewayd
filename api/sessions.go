package api

import (
	"encoding/json"
	"log"
	"net/http"

	"gatewayd/backend/global"
	"gatewayd/backend/session"

	"gatewayd/pkg/encoder"
	"github.com/go-martini/martini"
)

// API talks to session manager to get sessions.
// Remember sessions live each in it's own goroutine.

// SessionByToken responds to the API request to the session by token
func SessionByToken(params martini.Params, enc encoder.Encoder) (int, []byte) {
	token, ok := params["token"]
	if !ok {
		return http.StatusBadRequest, []byte("No token passed")
	}
	log.Printf("api: requesting session info for token %q", token)

	session, err := global.SessionManager.SessionByToken(token)
	if err != nil {
		log.Println(err)
		return http.StatusNotFound, []byte("No such session")
	}

	info := struct {
		State       string `json:"state"`
		ProfileName string `json:"profile"`
	}{
		session.Driver().State().String(),
		session.Profile().Name,
	}
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
	log.Println("api: creating session for request", t)

	profile, err := global.ProfileManager.Get(t.ProfileName)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, []byte("Unable to locate profile")
	}

	session, err := session.Create(profile)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, []byte("Unable to create session")
	}

	token, err := global.SessionManager.Manage(session)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, []byte("Unable to register session")
	}

	log.Printf("api: session created and registered (token %q)", token)

	go session.Run()

	result := struct {
		Token string `json:"token"`
	}{token}
	return http.StatusCreated, encoder.Must(enc.Encode(result))
}
