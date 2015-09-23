package gitpuller

import (
	"os"
	"strconv"
	"time"

	"github.com/carlosdp/harbor/options"
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

func (gp *gitPuller) Pull(url, commitHash string, ops options.Options) (string, error) {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	path := os.TempDir() + timeStamp
	repo, err := git.Clone(url, path, &git.CloneOptions{})
	if err != nil {
		return "", err
	}

	oid, err := git.NewOid(commitHash)
	if err != nil {
		return "", err
	}

	commit, err := repo.LookupCommit(oid)
	if err != nil {
		return "", err
	}

	tree, err := commit.Tree()
	if err != nil {
		return "", err
	}

	err = repo.CheckoutTree(tree, &git.CheckoutOpts{})
	if err != nil {
		return "", err
	}

	gp.workDir = repo.Workdir()

	return repo.Workdir(), nil
}
