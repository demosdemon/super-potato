package platformsh

type Mount struct {
	Source     ApplicationMount `json:"source"`
	SourcePath string           `json:"path"`
	Service    string           `json:"service,omitempty"`
}
