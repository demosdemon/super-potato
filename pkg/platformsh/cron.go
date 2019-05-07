package platformsh

type Cron struct {
	Spec string `json:"spec"`
	Cmd  string `json:"cmd"`
}
