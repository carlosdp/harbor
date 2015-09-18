package githubhook

import (
	"encoding/json"
	"net/http"
	"path"

	"github.com/carlosdp/harbor/hook"
)

type githubHook struct {
	FullName   string
	RepoURL    string
	Username   string
	CommitHash string
	Branch     string
}

type githubRepository struct {
	FullName string `json:"full_name"`
	SSHURL   string `json:"ssh_url"`
}

type githubRequest struct {
	CommitID string           `json:"after"`
	Repo     githubRepository `json:"repository"`
}

func init() {
	hook.RegisterHook("github-hook", &githubHook{})
}

func (gh *githubHook) New() hook.Hook {
	return &githubHook{}
}

func (gh *githubHook) HandleRequest(req *http.Request) error {
	decoder := json.NewDecoder(req.Body)
	var r githubRequest
	err := decoder.Decode(&r)
	if err != nil {
		return err
	}

	gh.FullName = r.Repo.FullName
	gh.RepoURL = r.Repo.SSHURL
	gh.CommitHash = r.CommitID

	return nil
}

func (gh *githubHook) Name() string {
	return path.Base(gh.FullName)
}

func (gh *githubHook) DeploymentID() string {
	return gh.CommitHash
}

func (gh *githubHook) URI() string {
	return gh.RepoURL
}
