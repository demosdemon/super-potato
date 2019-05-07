package platformsh

type Preflight struct {
	Enabled      bool     `json:"enabled"`
	IgnoredRules []string `json:"ignored_rules"`
}
