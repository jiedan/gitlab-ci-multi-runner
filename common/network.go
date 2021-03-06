package common

type UpdateState int

const (
	UpdateSucceeded UpdateState = iota
	UpdateAbort
	UpdateFailed
)

type FeaturesInfo struct {
	Variables bool `json:"variables"`
	Image     bool `json:"image"`
	Services  bool `json:"services"`
	Artifacts bool `json:"features"`
	Cache     bool `json:"cache"`
}

type VersionInfo struct {
	Name         string       `json:"name,omitempty"`
	Version      string       `json:"version,omitempty"`
	Revision     string       `json:"revision,omitempty"`
	Platform     string       `json:"platform,omitempty"`
	Architecture string       `json:"architecture,omitempty"`
	Executor     string       `json:"executor,omitempty"`
	Features     FeaturesInfo `json:"features"`
}

type GetBuildRequest struct {
	Info  VersionInfo `json:"info,omitempty"`
	Token string      `json:"token,omitempty"`
}

type BuildOptions map[string]interface{}

type GetBuildResponse struct {
	ID            int            `json:"id,omitempty"`
	ProjectID     int            `json:"project_id,omitempty"`
	Commands      string         `json:"commands,omitempty"`
	RepoURL       string         `json:"repo_url,omitempty"`
	Sha           string         `json:"sha,omitempty"`
	RefName       string         `json:"ref,omitempty"`
	BeforeSha     string         `json:"before_sha,omitempty"`
	AllowGitFetch bool           `json:"allow_git_fetch,omitempty"`
	Timeout       int            `json:"timeout,omitempty"`
	Variables     BuildVariables `json:"variables"`
	Options       BuildOptions   `json:"options"`
	Token         string         `json:"token"`
	Name          string         `json:"name"`
	Stage         string         `json:"stage"`
	Tag           bool           `json:"tag"`
	TLSCAChain    string         `json:"-"`
}

type RegisterRunnerRequest struct {
	Info        VersionInfo `json:"info,omitempty"`
	Token       string      `json:"token,omitempty"`
	Description string      `json:"description,omitempty"`
	Tags        string      `json:"tag_list,omitempty"`
}

type RegisterRunnerResponse struct {
	Token string `json:"token,omitempty"`
}

type DeleteRunnerRequest struct {
	Token string `json:"token,omitempty"`
}

type VerifyRunnerRequest struct {
	Token string `json:"token,omitempty"`
}

type UpdateBuildRequest struct {
	Info  VersionInfo `json:"info,omitempty"`
	Token string      `json:"token,omitempty"`
	State BuildState  `json:"state,omitempty"`
	Trace string      `json:"trace,omitempty"`
}

type Network interface {
	GetBuild(config RunnerConfig) (*GetBuildResponse, bool)
	RegisterRunner(config RunnerCredentials, description, tags string) *RegisterRunnerResponse
	DeleteRunner(config RunnerCredentials) bool
	VerifyRunner(config RunnerCredentials) bool
	UpdateBuild(config RunnerConfig, id int, state BuildState, trace string) UpdateState
	GetArtifactsUploadURL(config RunnerCredentials, id int) string
}
