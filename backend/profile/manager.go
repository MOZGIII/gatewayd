package profile

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// Right now this is a simple JSON-config-based implementation.
// Tricky thing is we have to deal with sessions that have
// already been started and are referencing the profile:
// we should either wipe out all the sessions, or block
// any interaction with "locked" profile.
// For now we just don't allow any interaction with
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

// LoadJSON loads JSON files with a list of profiles.
func (manager *Manager) LoadJSON(filename string) error {
	reader, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer reader.Close()

	dec := json.NewDecoder(reader)

	var data struct {
		profiles []*Profile
	}
	if err := dec.Decode(&data); err != nil {
		return err
	}

	for _, profile := range data.profiles {
		if err := manager.register(profile); err != nil {
			return err
		}
	}

	return nil
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
