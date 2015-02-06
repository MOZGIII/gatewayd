package main

import (
	"log"

	"gatewayd/backend/control"
	"gatewayd/config"
	"gatewayd/server"
	"gatewayd/utils"
)

func main() {
	log.Println("Loading config...")
	cfg := config.LoadBuiltinConfig()

	go func() {
		log.Println("Initilializng singals catcher...")
		utils.CatchSignals(control.Control.Quit())
	}()

	go func() {
		log.Println("Staring servers...")
		server.StartAll(cfg)
	}()

	// Control routine as "main thread".
	control.Control.Run()

	log.Println("Done!")
}
