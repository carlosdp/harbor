package dockerbuilder

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/carlosdp/harbor/builder"
	"github.com/fsouza/go-dockerclient"
)

type dockerBuilder struct {
	imageName string
}

func init() {
	builder.RegisterBuilder("docker-builder", &dockerBuilder{})
}

func (db *dockerBuilder) New() builder.Builder {
	return &dockerBuilder{}
}

func (db *dockerBuilder) Build(workDir, image string) (string, error) {
	image, err := createImage(workDir, image)
	if err != nil {
		return "", err
	}

	db.imageName = image

	return image, nil
}

func createImage(workDir, originalImage string) (string, error) {
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
	tarPath := os.TempDir() + timeStamp + ".tar"

	fw, err := os.Create(tarPath)
	if err != nil {
		return "", err
	}
	defer fw.Close()

	tw := tar.NewWriter(fw)
	defer tw.Close()

	err = writeTarDirectory(workDir[:len(workDir)-1], "", tw)

	if err != nil {
		return "", err
	}

	imageName, err := createDockerImage(originalImage, tarPath)

	return imageName, nil
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

func createDockerImage(originalName, tarPath string) (string, error) {
	dockerHost := os.Getenv("DOCKER_HOST")

	if dockerHost == "" {
		dockerHost = "unix:///var/run/docker.sock"
	}

	client, err := docker.NewClient(dockerHost)

	if err != nil {
		return "", err
	}

	f, err := os.Open(tarPath)

	if err != nil {
		return "", err
	}

	imageName := originalName + "-" + strconv.FormatInt(time.Now().Unix(), 10)

	err = client.BuildImage(docker.BuildImageOptions{
		Name:         imageName,
		InputStream:  f,
		OutputStream: os.Stdout,
	})
	if err != nil {
		return "", err
	}

	return imageName, nil
}
