package platformsh

type Build struct {
	Flavor string `json:"flavor"`
	Caches Caches `json:"caches"`
}
