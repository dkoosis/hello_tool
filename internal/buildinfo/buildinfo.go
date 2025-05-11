package buildinfo

// These variables are populated by LDFLAGS during the build process
var (
	Version    string = "dev"
	CommitHash string = "unknown"
	BuildDate  string = "unknown"
)
