package main

import (
	"archive/tar"
	"fmt"
	"github.com/fsouza/go-dockerclient"
	"github.com/libgit2/git2go"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(response http.ResponseWriter, request *http.Request) {
		hook := &GithubHook{}
		hook.HandleRequest(request)
		workDir := gitPuller(hook.URI())
		dockerBuild(hook.Name(), workDir)
		dockerScheduler(hook.Name(), hook.DeploymentID())
	})

	server := &http.Server{
		Addr:    ":3002",
		Handler: mux,
	}

	fmt.Println("starting server")
	server.ListenAndServe()
}

func gitPuller(url string) string {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	path := os.TempDir() + timeStamp
	repo, err := git.Clone(url, path, &git.CloneOptions{})
	if err != nil {
		panic(err)
	}

	return repo.Workdir()
}

func dockerBuild(name string, workDir string) {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	tarPath := os.TempDir() + timeStamp + ".tar"

	fw, err := os.Create(tarPath)
	if err != nil {
		panic(err)
	}
	defer fw.Close()

	tw := tar.NewWriter(fw)
	defer tw.Close()

	err = writeTarDirectory(workDir[:len(workDir)-1], "", tw)

	if err != nil {
		panic(err)
	}

	dockerHost := os.Getenv("DOCKER_HOST")

	if dockerHost == "" {
		dockerHost = "unix:///var/run/docker.sock"
	}

	client, err := docker.NewClient(dockerHost)

	if err != nil {
		panic(err)
	}

	f, err := os.Open(tarPath)

	if err != nil {
		panic(err)
	}

	err = client.BuildImage(docker.BuildImageOptions{
		Name:         name,
		InputStream:  f,
		OutputStream: os.Stdout,
	})
	if err != nil {
		panic(err)
	}
}

func dockerScheduler(name string, id string) {
	containerName := name + "-" + id
	dockerHost := os.Getenv("DOCKER_HOST")

	if dockerHost == "" {
		dockerHost = "unix:///var/run/docker.sock"
	}

	client, err := docker.NewClient(dockerHost)

	if err != nil {
		panic(err)
	}

	containers, err := client.ListContainers(docker.ListContainersOptions{
		All: true,
	})
	if err != nil {
		panic(err)
	}

	for _, c := range containers {
		err = client.StopContainer(c.ID, 5)
		if err != nil {
			fmt.Println(err)
		}
		err = client.RemoveContainer(docker.RemoveContainerOptions{
			ID:            c.ID,
			RemoveVolumes: true,
		})
		if err != nil {
			panic(err)
		}
	}

	c, err := client.CreateContainer(docker.CreateContainerOptions{
		Name: containerName,
		Config: &docker.Config{
			AttachStdin:  true,
			AttachStdout: true,
			AttachStderr: true,
			Image:        name,
			Cmd:          []string{"/bin/echo", "Hail Voltron"},
		},
	})

	if err != nil {
		panic(err)
	}

	err = client.StartContainer(c.ID, &docker.HostConfig{})

	if err != nil {
		panic(err)
	}

	err = client.AttachToContainer(docker.AttachToContainerOptions{
		Container:    c.ID,
		OutputStream: os.Stdout,
		ErrorStream:  os.Stderr,
		Logs:         true,
		Stream:       false,
		Stdout:       true,
		Stderr:       true,
	})

	if err != nil {
		panic(err)
	}
}

func writeTarDirectory(path string, shortPath string, tw *tar.Writer) error {
	dir, err := os.Open(path)
	if err != nil {
		return err
	}

	defer dir.Close()

	fis, err := dir.Readdir(0)

	if err != nil {
		return err
	}

	for _, finfo := range fis {
		curPath := path + "/" + finfo.Name()

		fi, err := os.Stat(curPath)
		if err != nil {
			return err
		}

		if fi.IsDir() {
			writeTarDirectory(curPath, shortPath+"/"+finfo.Name(), tw)
		} else {
			filePath, err := filepath.EvalSymlinks(curPath)
			if err != nil {
				return err
			}

			fr, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer fr.Close()

			h := &tar.Header{
				Name:    finfo.Name(),
				Size:    fi.Size(),
				Mode:    int64(fi.Mode()),
				ModTime: fi.ModTime(),
			}

			err = tw.WriteHeader(h)
			if err != nil {
				return err
			}

			_, err = io.Copy(tw, fr)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
