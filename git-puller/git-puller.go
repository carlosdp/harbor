package gitpuller

import (
	"os"
	"strconv"
	"time"

	"github.com/carlosdp/harbor/puller"
	git "gopkg.in/libgit2/git2go.v22"
)

type GitPuller struct {
	workDir string
}

func init() {
	puller.RegisterPuller("git-puller", &GitPuller{})
}

func (gp *GitPuller) New() puller.Puller {
	return &GitPuller{}
}

func (gp *GitPuller) Pull(url string) (string, error) {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	path := os.TempDir() + timeStamp
	repo, err := git.Clone(url, path, &git.CloneOptions{})
	if err != nil {
		return "", err
	}

	gp.workDir = repo.Workdir()

	return repo.Workdir(), nil
}
