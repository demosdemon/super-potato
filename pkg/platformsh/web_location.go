package platformsh

type WebLocation struct {
	Root     string    `json:"root"`
	Expires  Duration  `json:"expires"`
	Passthru Passthru  `json:"passthru"`
	Scripts  bool      `json:"scripts"`
	Index    []string  `json:"index"`
	Allow    bool      `json:"allow"`
	Headers  StringMap `json:"headers"`
	Rules    WebRules  `json:"rules"`
}
