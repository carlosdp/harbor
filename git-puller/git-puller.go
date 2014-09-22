package gitpuller

import (
	"github.com/libgit2/git2go"
	"os"
	"strconv"
	"time"
)

func Pull(url string) string {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	path := os.TempDir() + timeStamp
	repo, err := git.Clone(url, path, &git.CloneOptions{})
	if err != nil {
		panic(err)
	}

	return repo.Workdir()
}
