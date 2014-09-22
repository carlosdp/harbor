package main

import (
	"encoding/json"
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

func (gh *GithubHook) Endpoint() string {
	return "/github"
}

func (gh *GithubHook) HandleRequest(req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var r githubRequest
	err := decoder.Decode(&r)
	if err != nil {
		panic(err)
	}

	gh.FullName = r.Repo.FullName
	gh.RepoURL = r.Repo.SshUrl
	gh.CommitHash = r.CommitID
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
