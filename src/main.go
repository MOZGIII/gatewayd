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
)

func checkArgs() {
	if *configFile == "" || *profilesFile == "" {
		flag.Usage()
		os.Exit(2)
	}
}

func main() {
	flag.Parse()
	checkArgs()

	log.Printf("main: registered drivers: %v", driver.Registry())

	log.Println("main: loading config...")
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("main: unable to load config: %s", err)
	}

	log.Println("main: loading profiles...")
	if err := global.ProfileManager.LoadJSONFile(*profilesFile); err != nil {
		log.Fatalf("main: unable to load profiles: %s", err)
	}
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
