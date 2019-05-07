package platformsh

import "encoding/json"

type Passthru struct {
	Enabled bool
	Path    string
}

func (p *Passthru) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &p.Path); err == nil {
		p.Enabled = true
		return nil
	}

	return json.Unmarshal(data, &p.Enabled)
}

func (p Passthru) MarshalJSON() ([]byte, error) {
	if p.Path == "" {
		return json.Marshal(p.Enabled)
	}

	return json.Marshal(p.Path)
}
