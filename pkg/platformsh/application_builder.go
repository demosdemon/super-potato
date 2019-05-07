package platformsh

type ApplicationBuilder struct {
	ApplicationCore
	Dependencies JSONObject `json:"dependencies"`
	Build        Build      `json:"build"`
	Source       Source     `json:"source"`
}
