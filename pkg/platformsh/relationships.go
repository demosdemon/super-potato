package platformsh

type (
	Relationships map[string][]Relationship

	Relationship struct {
		Cluster  string     `json:"cluster"`
		Fragment string     `json:"fragment"`
		Host     string     `json:"host"`
		Hostname string     `json:"hostname"`
		IP       string     `json:"ip"`
		Password string     `json:"password"`
		Path     string     `json:"path"`
		Port     int        `json:"port"`
		Public   bool       `json:"public"`
		Query    JSONObject `json:"query"`
		Rel      string     `json:"rel"`
		Scheme   string     `json:"scheme"`
		Service  string     `json:"service"`
		SSL      JSONObject `json:"ssl"`
		Type     string     `json:"type"`
		Username string     `json:"username"`
	}
)
