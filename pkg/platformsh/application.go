package platformsh

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TODO: replace repetitive code with Genny generics

const (
	AccessLevelViewer AccessLevel = iota
	AccessLevelContributor
	AccessLevelAdmin
	totalAccessLevels
)

const (
	AccessTypeSSH AccessType = iota
	totalAccessTypes
)

const (
	ApplicationMountLocal ApplicationMountType = iota
	ApplicationMountTemp
	ApplicationMountService
	totalApplicationMountTypes
)

const (
	ServiceSizeAuto ServiceSize = iota
	ServiceSizeSmall
	ServiceSizeMedium
	ServiceSizeLarge
	ServiceSizeExtraLarge
	ServiceSizeDoubleExtraLarge
	ServiceSizeQuadrupleExtraLarge
	totalServiceSizes
)

const (
	SocketFamilyTCP SocketFamily = iota
	SocketFamilyUnix
	totalSocketFamilies
)

const (
	SocketProtocolHTTP SocketProtocol = iota
	SocketProtocolFastCGI
	SocketProtocolUWSGI
	totalSocketProtocols
)

var (
	accessLevels = [totalAccessLevels]string{
		"viewer",
		"contributor",
		"admin",
	}

	accessTypes = [totalAccessTypes]string{
		"ssh",
	}

	applicationMountTypes = [totalApplicationMountTypes]string{
		"local",
		"tmp",
		"service",
	}

	serviceSizes = [totalServiceSizes]string{
		"AUTO",
		"S",
		"M",
		"L",
		"XL",
		"2XL",
		"4XL",
	}

	socketFamilies = [totalSocketFamilies]string{
		"tcp",
		"unix",
	}

	socketProtocols = [totalSocketProtocols]string{
		"http",
		"fastcgi",
		"uwsgi",
	}

	accessLevelMap = map[string]AccessLevel{
		"viewer":      AccessLevelViewer,
		"contributor": AccessLevelContributor,
		"admin":       AccessLevelAdmin,
	}

	accessTypeMap = map[string]AccessType{
		"ssh": AccessTypeSSH,
	}

	applicationMountTypeMap = map[string]ApplicationMountType{
		"local":   ApplicationMountLocal,
		"tmp":     ApplicationMountTemp,
		"service": ApplicationMountService,
	}

	serviceSizeMap = map[string]ServiceSize{
		"AUTO": ServiceSizeAuto,
		"S":    ServiceSizeSmall,
		"M":    ServiceSizeMedium,
		"L":    ServiceSizeLarge,
		"XL":   ServiceSizeExtraLarge,
		"2XL":  ServiceSizeDoubleExtraLarge,
		"4XL":  ServiceSizeQuadrupleExtraLarge,
	}

	socketFamilyMap = map[string]SocketFamily{
		"tcp":  SocketFamilyTCP,
		"unix": SocketFamilyUnix,
	}

	socketProtocolMap = map[string]SocketProtocol{
		"http":    SocketProtocolHTTP,
		"fastcgi": SocketProtocolFastCGI,
		"uwsgi":   SocketProtocolUWSGI,
	}
)

type (
	AccessLevel          uint8
	AccessType           uint8
	ApplicationMountType uint8
	ServiceSize          uint8
	SocketFamily         uint8
	SocketProtocol       uint8

	ApplicationAccess           map[AccessType]AccessLevel
	ApplicationCrons            map[string]ApplicationCron
	ApplicationMounts           map[string]ApplicationMount
	ApplicationRelationships    map[string]string
	ApplicationSourceOperations map[string]ApplicationSourceOperation
	ApplicationWorkers          map[string]ApplicationWorker
	Caches                      map[string]CacheConfiguration
	VariableNamespace           map[string]interface{}
	Variables                   map[string]VariableNamespace
	WebLocations                map[string]WebLocation
	WebRules                    map[string]WebRule
	ApplicationDependencies     map[string]interface{}

	ApplicationMount struct {
		Source     ApplicationMountType `json:"source"`
		SourcePath string               `json:"path"`
		Service    string               `json:"service,omitempty"`
	}

	BaseApplication struct {
		Size          ServiceSize              `json:"size"`
		Disk          uint32                   `json:"disk"`
		Access        ApplicationAccess        `json:"access"`
		Relationships ApplicationRelationships `json:"relationships"`
		Mounts        ApplicationMounts        `json:"mounts"`
		Timezone      string                   `json:"timezone"` // TODO: replace with time.Location
		Variables     Variables                `json:"variables"`
	}

	Passthru struct {
		Enabled bool
		Path    string
	}

	WebRule struct {
		Expires  Duration          `json:"expires"`
		Passthru Passthru          `json:"passthru"`
		Scripts  bool              `json:"scripts"`
		Allow    bool              `json:"allow"`
		Headers  map[string]string `json:"headers"`
	}

	WebLocation struct {
		Root     string            `json:"root"`
		Expires  Duration          `json:"expires"`
		Passthru Passthru          `json:"passthru"`
		Scripts  bool              `json:"scripts"`
		Index    []string          `json:"index"`
		Allow    bool              `json:"allow"`
		Headers  map[string]string `json:"headers"`
		Rules    WebRules          `json:"rules"`
	}

	Commands struct {
		Start string `json:"start"`
		Stop  string `json:"stop,omitempty"`
	}

	Upstream struct {
		SocketFamily SocketFamily   `json:"socket_family"`
		Protocol     SocketProtocol `json:"socket_protocol"`
	}

	WebApplication struct {
		// BaseApplication
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

	CacheConfiguration struct {
		Directory        string   `json:"directory"`
		Watch            []string `json:"watch"`
		AllowStale       bool     `json:"allow_stale"`
		ShareBetweenApps bool     `json:"share_between_apps"`
	}

	ApplicationBuild struct {
		Flavor string `json:"flavor"`
		Caches Caches `json:"caches"`
	}

	ApplicationSourceOperation struct {
		Command string `json:"command"`
	}

	ApplicationSource struct {
		Operations ApplicationSourceOperations `json:"operations"`
	}

	ApplicationHooks struct {
		Build      string `json:"build"`
		Deploy     string `json:"deploy"`
		PostDeploy string `json:"post_deploy"`
	}

	ApplicationPreflight struct {
		Enabled      bool     `json:"enabled"`
		IgnoredRules []string `json:"ignored_rules"`
	}

	ApplicationCron struct {
		Spec string `json:"spec"`
		Cmd  string `json:"cmd"`
	}

	ApplicationWorker struct {
		// BaseApplication
		Commands Commands `json:"commands"`
	}

	ApplicationCore struct {
		BaseApplication
		Name      string               `json:"name"`
		Type      string               `json:"type"`
		Runtime   interface{}          `json:"runtime"`
		Preflight ApplicationPreflight `json:"preflight"`
	}

	ApplicationBuilder struct {
		ApplicationCore
		Dependencies ApplicationDependencies `json:"dependencies"`
		Build        ApplicationBuild        `json:"build"`
		Source       ApplicationSource       `json:"source"`
	}

	Application struct {
		ApplicationCore
		Web     WebApplication     `json:"web"`
		Hooks   ApplicationHooks   `json:"hooks"`
		Crons   ApplicationCrons   `json:"crons"`
		Workers ApplicationWorkers `json:"workers"`
		TreeID  string             `json:"tree_id"`
		SlugID  string             `json:"slug_id"`
		AppDir  string             `json:"app_dir"`
	}
)

func NewAccessLevel(name string) (AccessLevel, error) {
	if v, ok := accessLevelMap[name]; ok {
		return v, nil
	}

	return 0, fmt.Errorf("unknown access level %q", name)
}

func (l AccessLevel) String() string {
	if l < totalAccessLevels {
		return accessLevels[l]
	}
	return fmt.Sprintf("unknown access level %02x", uint8(l))
}

func (l *AccessLevel) UnmarshalText(text []byte) (err error) {
	*l, err = NewAccessLevel(string(text))
	return err
}

func (l AccessLevel) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func NewAccessType(name string) (AccessType, error) {
	if v, ok := accessTypeMap[name]; ok {
		return v, nil
	}

	return 0, fmt.Errorf("unknown access type %q", name)
}

func (t AccessType) String() string {
	if t < totalAccessTypes {
		return accessTypes[t]
	}
	return fmt.Sprintf("unknown access type %02x", uint8(t))
}

func (t *AccessType) UnmarshalText(text []byte) (err error) {
	*t, err = NewAccessType(string(text))
	return err
}

func (t AccessType) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func NewApplicationMountType(name string) (ApplicationMountType, error) {
	if v, ok := applicationMountTypeMap[name]; ok {
		return v, nil
	}

	return 0, fmt.Errorf("unknown application mount type %q", name)
}

func (t ApplicationMountType) String() string {
	if t < totalApplicationMountTypes {
		return applicationMountTypes[t]
	}
	return fmt.Sprintf("unknown application mount type %02x", uint8(t))
}

func (t *ApplicationMountType) UnmarshalText(text []byte) (err error) {
	*t, err = NewApplicationMountType(string(text))
	return err
}

func (t ApplicationMountType) MarshalText() ([]byte, error) {
	return []byte(t.String()), nil
}

func NewServiceSize(name string) (ServiceSize, error) {
	if v, ok := serviceSizeMap[name]; ok {
		return v, nil
	}

	return 0, fmt.Errorf("unknown service size %q", name)
}

func (s ServiceSize) String() string {
	if s < totalServiceSizes {
		return serviceSizes[s]
	}
	return fmt.Sprintf("unknown service size %02x", uint8(s))
}

func (s *ServiceSize) UnmarshalText(text []byte) (err error) {
	*s, err = NewServiceSize(string(text))
	return err
}

func (s ServiceSize) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (t *ApplicationMount) UnmarshalText(text []byte) error {
	s := string(text)
	if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
		return json.Unmarshal(text, t)
	}

	if strings.HasPrefix(s, "shared:files") {
		t.Source = ApplicationMountLocal
		t.SourcePath = strings.TrimPrefix(s, "shared:files")
		return nil
	}

	return fmt.Errorf("invalid application mount %q", s)
}

func NewSocketFamily(name string) (SocketFamily, error) {
	if v, ok := socketFamilyMap[name]; ok {
		return v, nil
	}

	return 0, fmt.Errorf("unknown socket family %q", name)
}

func (f SocketFamily) String() string {
	if f < totalSocketFamilies {
		return socketFamilies[f]
	}

	return fmt.Sprintf("unknown socket family %02x", uint8(f))
}

func (f *SocketFamily) UnmarshalText(text []byte) (err error) {
	*f, err = NewSocketFamily(string(text))
	return err
}

func (f SocketFamily) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

func NewSocketProtocol(name string) (SocketProtocol, error) {
	if v, ok := socketProtocolMap[name]; ok {
		return v, nil
	}
	return 0, fmt.Errorf("unknown socket protocol %q", name)
}

func (p SocketProtocol) String() string {
	if p < totalSocketProtocols {
		return socketProtocols[p]
	}
	return fmt.Sprintf("unknown sockect protocol %02x", uint8(p))
}

func (p *SocketProtocol) UnmarshalText(text []byte) (err error) {
	*p, err = NewSocketProtocol(string(text))
	return err
}

func (p SocketProtocol) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
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
