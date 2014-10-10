package githubhook

import (
	"encoding/json"
	"github.com/carlosdp/harbor/hook"
	"net/http"
	"path"
)

type GithubHook struct {
	FullName   string
	RepoURL    string
	Username   string
	CommitHash string
	Branch     string
}

type githubRepository struct {
	FullName string `json:"full_name"`
	SshUrl   string `json:"ssh_url"`
}

type githubRequest struct {
	CommitID string           `json:"after"`
	Repo     githubRepository `json:"repository"`
}

func init() {
	hook.RegisterHook("github-hook", &GithubHook{})
}

func (gh *GithubHook) New() hook.Hook {
	return &GithubHook{}
}

func (gh *GithubHook) HandleRequest(req *http.Request) error {
	decoder := json.NewDecoder(req.Body)
	var r githubRequest
	err := decoder.Decode(&r)
	if err != nil {
		return err
	}

	gh.FullName = r.Repo.FullName
	gh.RepoURL = r.Repo.SshUrl
	gh.CommitHash = r.CommitID

	return nil
}

func (gh *GithubHook) Name() string {
	return path.Base(gh.FullName)
}

func (gh *GithubHook) DeploymentID() string {
	return gh.CommitHash
}

func (gh *GithubHook) URI() string {
	return gh.RepoURL
}
