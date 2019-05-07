package platformsh

type ApplicationCore struct {
	ApplicationBase
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	Runtime   interface{} `json:"runtime"`
	Preflight Preflight   `json:"preflight"`
}
