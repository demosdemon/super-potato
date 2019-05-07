package platformsh

type Commands struct {
	Start string `json:"start"`
	Stop  string `json:"stop,omitempty"`
}
