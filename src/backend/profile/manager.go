package profile

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// Right now this is a simple JSON-config-based implementation.
// Tricky thing is we have to deal with sessions that have
// already been started and are referencing the profile:
// we should either wipe out all the sessions, or block
// any interaction with "locked" profile.
// For now we just don't support any interaction with
// profiles after they are loaded.

// Manager stores profiles.
type Manager struct {
	profileByName map[string]*Profile
	mu            sync.RWMutex
}

// NewManager creates new profile manager.
func NewManager() *Manager {
	return &Manager{
		make(map[string]*Profile),
		sync.RWMutex{},
	}
}

func (manager *Manager) register(profile *Profile) error {
	if profile == nil {
		return fmt.Errorf("profilemanager: unable to register nil profile")
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()

	if _, exists := manager.profileByName[profile.Name]; exists {
		return fmt.Errorf("profilemanager: profile with name %q already registered", profile.Name)
	}

	manager.profileByName[profile.Name] = profile
	return nil
}

// LoadJSON loads profiles in JSON format from a reader.
func (manager *Manager) LoadJSON(reader io.Reader) error {
	dec := json.NewDecoder(reader)

	var data struct {
		Profiles []*Profile `json:"profiles"`
	}
	if err := dec.Decode(&data); err != nil {
		return err
	}

	for _, profile := range data.Profiles {
		if err := manager.register(profile); err != nil {
			return err
		}
	}

	return nil
}

// LoadJSONFile loads profiles from a JSON file.
func (manager *Manager) LoadJSONFile(filename string) error {
	reader, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer reader.Close()
	return manager.LoadJSON(reader)
}

// Get returns profile for given name.
func (manager *Manager) Get(name string) (*Profile, error) {
	manager.mu.RLock()
	defer manager.mu.RUnlock()

	profile, present := manager.profileByName[name]
	if !present {
		return nil, fmt.Errorf("profilemanager: profile with name %q not found", name)
	}

	return profile, nil
}

// Report dumps info about the current state
func (manager *Manager) Report() {
	log.Println("profilemanager: report", manager.profileByName)
}
