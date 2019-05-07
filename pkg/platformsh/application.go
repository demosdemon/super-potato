package platformsh

type Application struct {
	ApplicationCore
	Web     Web     `json:"web"`
	Hooks   Hooks   `json:"hooks"`
	Crons   Crons   `json:"crons"`
	Workers Workers `json:"workers"`
	TreeID  string  `json:"tree_id"`
	SlugID  string  `json:"slug_id"`
	AppDir  string  `json:"app_dir"`
}
