package platformsh

type CacheConfiguration struct {
	Directory        string   `json:"directory"`
	Watch            []string `json:"watch"`
	AllowStale       bool     `json:"allow_stale"`
	ShareBetweenApps bool     `json:"share_between_apps"`
}
