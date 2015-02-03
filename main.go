package main

import (
	"log"
	"net/http"

	"gatewayd/backend"
	"gatewayd/config"
	"gatewayd/server"
	"gatewayd/utils"
)

import _ "expvar"

func main() {
	log.Println("Loading config...")
	cfg := config.LoadBuiltinConfig()

	go func() {
		log.Println("Staring servers...")
		server.StartAll(cfg)
	}()

	go func() {
		log.Println("Initilializng singals catcher...")
		utils.CatchSignals(backend.Control.Quit())
	}()

	go func() {
		log.Println("Starting expvar...")
		http.ListenAndServe(":3010", nil)
	}()

	backend.Control.Run()

	log.Println("Done!")
}
