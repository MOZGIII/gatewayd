package main

import (
	"flag"
	// "fmt"
	"log"
	"os"

	"gatewayd/driver"

	"gatewayd/backend/control"
	"gatewayd/backend/global"
	"gatewayd/config"
	"gatewayd/server"
	"gatewayd/utils"

	_ "gatewayd/driver/localexec"
)

var (
	configFile   = flag.String("config", "", "path to config.json")
	profilesFile = flag.String("profiles", "", "path to profiles.json")

	currentConfig *config.Config
)

func checkArgs() {
	if *configFile == "" || *profilesFile == "" {
		flag.Usage()
		os.Exit(2)
	}
}

func loadConfiguration() {
	log.Println("main: loading config...")
	var err error
	currentConfig, err = config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("main: unable to load config: %s", err)
	}

	log.Println("main: loading profiles...")
	if err := global.ProfileManager.LoadJSONFile(*profilesFile); err != nil {
		log.Fatalf("main: unable to load profiles: %s", err)
	}
}

func main() {
	flag.Parse()
	checkArgs()

	// Print available drivers
	log.Printf("main: registered drivers: %v", driver.Registry())

	// Load external configuration files
	loadConfiguration()

	// Report on loaded profiles
	global.ProfileManager.Report()

	go func() {
		log.Println("main: starting singals catcher...")
		utils.CatchSignals(control.Control.Quit())
	}()

	go func() {
		log.Println("main: starting servers...")
		server.StartAll(currentConfig)
	}()

	// Control routine as "main thread".
	log.Println("main: starting control...")
	control.Control.Run()

	log.Println("Done!")
}
