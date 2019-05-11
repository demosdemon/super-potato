package platformsh

type (
	JSONArray  = []interface{}
	JSONObject = map[string]interface{}
	StringMap  = map[string]string

	Access           map[AccessType]AccessLevel
	Caches           map[string]CacheConfiguration
	Crons            map[string]Cron
	Mounts           map[string]Mount
	SourceOperations map[string]SourceOperation
	Variables        map[string]JSONObject
	WebLocations     map[string]WebLocation
	WebRules         map[string]WebRule
	Workers          map[string]Worker

	Application struct {
		ApplicationCore
		Web     Web     `json:"web"`
		Hooks   Hooks   `json:"hooks"`
		Crons   Crons   `json:"crons"`
		Workers Workers `json:"workers"`
		TreeID  string  `json:"tree_id"`
		SlugID  string  `json:"slug_id"`
		AppDir  string  `json:"app_dir"`
	}

	ApplicationBase struct {
		Size          ServiceSize `json:"size"`
		Disk          uint32      `json:"disk"`
		Access        Access      `json:"access"`
		Relationships StringMap   `json:"relationships"`
		Mounts        Mounts      `json:"mounts"`
		Timezone      string      `json:"timezone"` // TODO: replace with serializable time.Location
		Variables     Variables   `json:"variables"`
	}

	ApplicationBuilder struct {
		ApplicationCore
		Dependencies JSONObject `json:"dependencies"`
		Build        Build      `json:"build"`
		Source       Source     `json:"source"`
	}

	ApplicationCore struct {
		ApplicationBase
		Name      string      `json:"name"`
		Type      string      `json:"type"`
		Runtime   interface{} `json:"runtime"`
		Preflight Preflight   `json:"preflight"`
	}

	Build struct {
		Flavor string `json:"flavor"`
		Caches Caches `json:"caches"`
	}

	CacheConfiguration struct {
		Directory        string   `json:"directory"`
		Watch            []string `json:"watch"`
		AllowStale       bool     `json:"allow_stale"`
		ShareBetweenApps bool     `json:"share_between_apps"`
	}

	Commands struct {
		Start string `json:"start"`
		Stop  string `json:"stop,omitempty"`
	}

	Cron struct {
		Spec string `json:"spec"`
		Cmd  string `json:"cmd"`
	}

	Hooks struct {
		Build      string `json:"build"`
		Deploy     string `json:"deploy"`
		PostDeploy string `json:"post_deploy"`
	}

	Mount struct {
		Source     ApplicationMount `json:"source"`
		SourcePath string           `json:"path"`
		Service    string           `json:"service,omitempty"`
	}

	Preflight struct {
		Enabled      bool     `json:"enabled"`
		IgnoredRules []string `json:"ignored_rules"`
	}

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

	Source struct {
		Operations SourceOperations `json:"operations"`
	}

	SourceOperation struct {
		Command string `json:"command"`
	}

	Upstream struct {
		SocketFamily SocketFamily   `json:"socket_family"`
		Protocol     SocketProtocol `json:"socket_protocol"`
	}

	Web struct {
		// ApplicationBase
		Locations    WebLocations `json:"locations"`
		Commands     Commands     `json:"commands"`
		Upstream     Upstream     `json:"upstream"`
		DocumentRoot *string      `json:"document_root,omitempty"` // deprecated
		Passthru     *string      `json:"passthru,omitempty"`      // deprecated
		IndexFiles   []string     `json:"index_files,omitempty"`   // deprecated
		Whitelist    []string     `json:"whitelist,omitempty"`     // deprecated
		Blacklist    []string     `json:"blacklist,omitempty"`     // deprecated
		Expires      *Duration    `json:"expires,omitempty"`       // deprecated
		MoveToRoot   *bool        `json:"move_to_root,omitempty"`  // deprecated
	}

	WebLocation struct {
		Root     string    `json:"root"`
		Expires  Duration  `json:"expires"`
		Passthru Passthru  `json:"passthru"`
		Scripts  bool      `json:"scripts"`
		Index    []string  `json:"index"`
		Allow    bool      `json:"allow"`
		Headers  StringMap `json:"headers"`
		Rules    WebRules  `json:"rules"`
	}

	WebRule struct {
		Expires  Duration  `json:"expires"`
		Passthru Passthru  `json:"passthru"`
		Scripts  bool      `json:"scripts"`
		Allow    bool      `json:"allow"`
		Headers  StringMap `json:"headers"`
	}

	Worker struct {
		// ApplicationBase
		Commands Commands `json:"commands"`
	}
)
