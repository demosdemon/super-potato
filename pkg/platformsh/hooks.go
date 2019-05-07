package platformsh

type Hooks struct {
	Build      string `json:"build"`
	Deploy     string `json:"deploy"`
	PostDeploy string `json:"post_deploy"`
}
