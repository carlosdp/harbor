package githubhook

import (
	"encoding/json"
	"net/http"
	"path"

	log "github.com/carlosdp/supply-chain/Godeps/_workspace/src/github.com/Sirupsen/logrus"
	"github.com/carlosdp/supply-chain/hook"
	"github.com/carlosdp/supply-chain/options"
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

func (gh *githubHook) Start(mux *http.ServeMux, queue chan<- hook.Request, ops options.Options) error {
	mux.HandleFunc(ops.GetString("endpoint"), func(res http.ResponseWriter, req *http.Request) {
		decoder := json.NewDecoder(req.Body)
		var r githubRequest
		err := decoder.Decode(&r)
		if err != nil {
			log.Error(err)
		}

		res.WriteHeader(200)

		request := hook.NewRequest(r.Repo.FullName, r.CommitID, r.Repo.SSHURL, r.Repo.FullName+"-"+r.CommitID)
		queue <- request
	})

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
