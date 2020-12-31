package cli

var (
	cli_name     = "kmgmt"
	cli_version  = "unknown"
	cli_gitsha   = "unknown sha"
	cli_gitdirty = ""
)

type CompiledEnv struct {
	Name     string
	Version  string
	GitSha   string
	GitDirty bool
}

var env CompiledEnv

func init() {
	// must be created inside the init function to pickup build specific params
	env = CompiledEnv{
		Name:     cli_name,
		Version:  cli_version,
		GitSha:   cli_gitsha,
		GitDirty: cli_gitdirty != "",
	}
}
