package docker

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

// DockerClient is a slimmed down implementation of the docker cli
type DockerClient struct {
	Host            string
	Cli             *client.Client
	HostConf        *container.HostConfig
	NetConf         *network.NetworkingConfig
	Conf            *container.Config
	Registry, Image string
	Entrypoint      []string
}

// NewClient returns a new DockerClient initialised with the API object
func NewClient() *DockerClient {
	c := new(DockerClient)
	c.SetDefaults()
	return c
}

// InitDocker initialises the client
func (c *DockerClient) InitDocker() error {
	var cli *client.Client

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	cli, err := client.NewClientWithOpts(
		client.WithHost(c.Host),
		client.WithVersion("v1.37"),
		client.WithHTTPHeaders(defaultHeaders))

	if err != nil {
		return err
	}
	c.Cli = cli
	return nil
}

// SetDefaults sets container, host and net configs to defaults. Called when instantiating a new client or can be called
// manually at any time to reset API configs back to empty defaults
func (c *DockerClient) SetDefaults() {
	c.HostConf = &container.HostConfig{Binds: []string{}}
	c.NetConf = &network.NetworkingConfig{}
	c.Conf = &container.Config{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		Tty:          true,
		Env:          []string{},
	}
}

// SetHostConf sets the container.HostConfig struct for the new container
func (c *DockerClient) SetHostConf(h *container.HostConfig) {
	c.HostConf = h
}

// SetNetConf sets the network.NetworkingConfig struct for the new container
func (c *DockerClient) SetNetConf(n *network.NetworkingConfig) {
	c.NetConf = n
}

// SetConf sets the container.Config struct for the new container
func (c *DockerClient) SetConf(co *container.Config) {
	c.Conf = co
}

// AddBind adds a bind mount to the HostConfig
func (c *DockerClient) AddBind(bnd string) {
	c.HostConf.Binds = append(c.HostConf.Binds, bnd)
}

// AddEnv adds an environment variable to the HostConfig
func (c *DockerClient) AddEnv(key, value string) {
	c.Conf.Env = append(c.Conf.Env, fmt.Sprintf("%s=%s", key, value))
}

// AddBinds adds multiple bind mounts to the HostConfig
func (c *DockerClient) AddBinds(bnds []string) {
	c.HostConf.Binds = append(c.HostConf.Binds, bnds...)
}

// AddEnvs adds multiple envs to the HostConfig
func (c *DockerClient) AddEnvs(envs []string) {
	c.Conf.Env = append(c.Conf.Env, envs...)
}

// SetBinds sets the bind mounts in the HostConfig
func (c *DockerClient) SetBinds(bnds []string) {
	c.HostConf.Binds = bnds
}

// SetEnvs sets the environment variables in the Conf
func (c *DockerClient) SetEnvs(envs []string) {
	c.Conf.Env = envs
}

// SetImage sets the image in Conf
func (c *DockerClient) SetImage(img string) {
	c.Image = img
	c.setImage()
}

// SetCmd sets the command the Entrypoint
func (c *DockerClient) SetEntrypoint(ep []string) {
	c.Entrypoint = ep
}

func (c *DockerClient) setImage() {
	var img string
	if c.Registry != "" {
		img = c.Registry + "/" + c.Image
	} else {
		img = c.Image
	}
	c.Conf.Image = img
}

// SetRegistry sets the registry where to pull images from
func (c *DockerClient) SetRegistry(reg string) {
	c.Registry = reg
	c.setImage()
}

// Privileged sets whether the container should run as privileged
func (c *DockerClient) Privileged(p bool) {
	c.HostConf.Privileged = p
}

// SetCmd sets the command to run in the container
func (c *DockerClient) SetCmd(cmd []string) {
	c.Conf.Cmd = cmd
}

// SetWorkDir sets the working directory of the container
func (c *DockerClient) SetWorkDir(wd string) {
	c.Conf.WorkingDir = wd
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
