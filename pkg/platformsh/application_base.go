package platformsh

type ApplicationBase struct {
	Size          ServiceSize `json:"size"`
	Disk          uint32      `json:"disk"`
	Access        Access      `json:"access"`
	Relationships StringMap   `json:"relationships"`
	Mounts        Mounts      `json:"mounts"`
	Timezone      string      `json:"timezone"` // TODO: replace with serializable time.Location
	Variables     Variables   `json:"variables"`
}
