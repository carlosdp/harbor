package gitpuller

import (
	"os"
	"strconv"
	"time"

	"github.com/carlosdp/harbor/puller"
	git "gopkg.in/libgit2/git2go.v22"
)

type gitPuller struct {
	workDir string
}

func init() {
	puller.RegisterPuller("git-puller", &gitPuller{})
}

func (gp *gitPuller) New() puller.Puller {
	return &gitPuller{}
}

func (gp *gitPuller) Pull(url string) (string, error) {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	path := os.TempDir() + timeStamp
	repo, err := git.Clone(url, path, &git.CloneOptions{})
	if err != nil {
		return "", err
	}

	gp.workDir = repo.Workdir()

	return repo.Workdir(), nil
}