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
)
