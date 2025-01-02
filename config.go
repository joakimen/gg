package clone

type Config struct {
	Owner           string
	Repo            string
	IncludeArchived bool
	Limit           int
	CloneDir        string
	Verbose         bool
	RepoFile        string
}

type Flags struct {
	Owner           string
	Repo            string
	IncludeArchived bool
	Limit           int
	CloneDir        string
	Verbose         bool
	RepoFile        string
}

type Envs struct {
	CloneDir string
}
