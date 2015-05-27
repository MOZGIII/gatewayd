package main

import (
	"log"

	"gatewayd/driver"

	"gatewayd/backend/control"
	"gatewayd/backend/global"
	_ "gatewayd/config"
	"gatewayd/server"
	"gatewayd/utils"

	_ "gatewayd/driver/localexec"

	"gatewayd/defaultconfigs"
)

func main() {
	log.Printf("main: registered drivers: %v", driver.Registry())

	log.Println("main: loading config...")
	cfg := defaultconfigs.GetConfig()

	log.Println("main: loading profiles...")
	global.ProfileManager.LoadJSON(defaultconfigs.GetProfileJSON())
	global.ProfileManager.Report()

	go func() {
		log.Println("main: starting singals catcher...")
		utils.CatchSignals(control.Control.Quit())
	}()

	go func() {
		log.Println("main: starting servers...")
		server.StartAll(cfg)
	}()

	// Control routine as "main thread".
	log.Println("main: starting control...")
	control.Control.Run()

	log.Println("Done!")
}
