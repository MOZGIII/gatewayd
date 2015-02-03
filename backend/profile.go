package backend

// Profile holds data passed by foreign users of our system
type Profile struct {
	Id          int
	ForeignKeys map[string]string // we will need to maintain index of these... big global nonpersistent key/value?
	Params      map[string]string // key/value pairs that are passed around for the system to use
}

// Get profile by id
func ProfileById(id int) (*Profile, error) {
	// get profile from db
}
