package platformsh

type WebRule struct {
	Expires  Duration  `json:"expires"`
	Passthru Passthru  `json:"passthru"`
	Scripts  bool      `json:"scripts"`
	Allow    bool      `json:"allow"`
	Headers  StringMap `json:"headers"`
}
