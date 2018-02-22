package cali

import (
	"crypto/md5"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

// GitCheckoutConfig is input for Git.Checkout
type GitCheckoutConfig struct {
	Repo, Branch, RelPath, Image string
}

const gitImage = "indiehosters/git:latest"

// Git returns a new instance
func (c *DockerClient) Git() *Git {
	return &Git{c: c, Image: gitImage}
}

// Git is used to interact with containerised git
type Git struct {
	c     *DockerClient
	Image string
}

// Checkout will create and start a container, checkout repo and leave container stopped
// so volume can be imported
func (g *Git) Checkout(cfg *GitCheckoutConfig) (string, error) {
	name := fmt.Sprintf("data_%x", md5.Sum([]byte(cfg.Repo+cfg.Branch)))

	if g.c.ContainerExists(name) {
		log.WithField("name", name).Info("Existing data container found")

		if _, err := g.Pull(name); err != nil {
			log.Warnf("Git pull error: %s", err)
			return name, err
		}
		return name, nil
	}

	log.WithFields(log.Fields{
		"git_url": cfg.Repo,
		"image":   g.Image,
	}).Info("Creating data containers")

	co := container.Config{
		Cmd:          []string{"clone", cfg.Repo, "-b", cfg.Branch, "--depth", "1", "."},
		Image:        gitImage,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   "/tmp/workspace",
		Entrypoint:   []string{"git"},
	}
	hc := container.HostConfig{
		Binds: []string{
			"/tmp/workspace",
			fmt.Sprintf("%s/.ssh:/root/.ssh", os.Getenv("HOME")),
		},
	}
	nc := network.NetworkingConfig{}

	g.c.SetConf(&co)
	g.c.SetHostConf(&hc)
	g.c.SetNetConf(&nc)

	id, err := g.c.StartContainer(false, name)

	if err != nil {
		return "", err
	}
	return id, nil
}

// Pull a git repository
func (g *Git) Pull(name string) (string, error) {
	co := container.Config{
		Cmd:          []string{"pull"},
		Image:        g.Image,
		Tty:          true,
		AttachStdout: true,
		AttachStderr: true,
		WorkingDir:   "/tmp/workspace",
		Entrypoint:   []string{"git"},
	}
	hc := container.HostConfig{
		VolumesFrom: []string{name},
		Binds: []string{
			fmt.Sprintf("%s/.ssh:/root/.ssh", os.Getenv("HOME")),
		},
	}
	nc := network.NetworkingConfig{}

	g.c.SetConf(&co)
	g.c.SetHostConf(&hc)
	g.c.SetNetConf(&nc)

	return g.c.StartContainer(true, "")
}
