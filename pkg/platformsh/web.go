package platformsh

type Web struct {
	// ApplicationBase
	Locations    WebLocations `json:"locations"`
	Commands     Commands     `json:"commands"`
	Upstream     Upstream     `json:"upstream"`
	DocumentRoot *string      `json:"document_root,omitempty"` // deprecated
	Passthru     *string      `json:"passthru,omitempty"`      // deprecated
	IndexFiles   []string     `json:"index_files,omitempty"`   // deprecated
	Whitelist    []string     `json:"whitelist,omitempty"`     // deprecated
	Blacklist    []string     `json:"blacklist,omitempty"`     // deprecated
	Expires      *Duration    `json:"expires,omitempty"`       // deprecated
	MoveToRoot   *bool        `json:"move_to_root,omitempty"`  // deprecated
}
