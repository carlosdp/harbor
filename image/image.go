package image

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/fsouza/go-dockerclient"
)

func NewImage(name string, workDir string) {
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
				Name:    shortPath + "/" + finfo.Name(),
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
