package main

import (
	"log"
	"strings"

	"gatewayd/driver"

	"gatewayd/backend/control"
	"gatewayd/backend/global"
	"gatewayd/config"
	"gatewayd/server"
	"gatewayd/utils"

	_ "gatewayd/driver/localexec"
)

func main() {
	log.Printf("main: registered drivers: %v", driver.Registry())

	log.Println("main: loading config...")
	cfg := config.LoadBuiltinConfig()

	log.Println("main: loading profiles...")
	global.ProfileManager.LoadJSON(strings.NewReader(`{
		"profiles": [
			{
				"name": "test",
				"driver": "localexec",
				"params": {
					"command": {
						"name": "gateway-session-test",
						"args": ["test"]
					}
				}
			}
		]
	}`))
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
