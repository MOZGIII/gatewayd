package main

import (
	"flag"
	"log"
	"os"

	"gatewayd/driver"
	"gatewayd/policy"

	"gatewayd/config"
	"gatewayd/global"
	"gatewayd/server"
	"gatewayd/utils"

	// Load drivers
	_ "gatewayd/driver/localexec"

	// Load policies
	_ "gatewayd/policy/sessionmanagement/noop"
	_ "gatewayd/policy/tunnelaccess/noop"
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

	// Print available drivers.
	log.Printf("main: registered drivers: %v", driver.Registry())

	// Load external configuration files.
	loadConfiguration()

	// Initialize policies with the loaded config.
	log.Println("main: initializing policies...")
	if err := policy.LoadFromConfig(&currentConfig.Policies); err != nil {
		log.Fatalf("main: unable to load policies: %s", err)
	}

	// Report on loaded profiles.
	global.ProfileManager.Report()

	log.Println("main: starting servers...")
	server.RunAll(currentConfig)

	// Quit channel is used to handle graceful shutdowns.
	quitch := make(chan struct{})

	log.Println("main: starting singals catcher...")
	go utils.CatchSignals(quitch)

	// Continue only if asked to exit.
	<-quitch

	// Make sure all runables quit gracefully.
	global.Runner.TerminateAll()
	global.Runner.Wait()

	log.Println("Done!")
}
