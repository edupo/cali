package docker

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/jhoonb/archivex"
	"golang.org/x/crypto/ssh/terminal"
	"golang.org/x/net/context"
)

// Containers are used by this version of cali as background processes.
// Cali containers usually run `sleep infinite` and make several exec's to
// perform the intended operations.
// Cali containers are also executed as the current user so the permissions
// are not messed-up. This is specially useful when running cali apps on CI.

// ExecContainer is the main functionality of this library. It initializes,
// fixes, executes and removes a container.
func (c *DockerClient) ExecContainer(rm bool, name string) (string, error) {

	// Runs the container just doing nothing
	c.Conf.Entrypoint = []string{"sleep", "infinity"}
	id, err := c.initializeContainer(name)
	if err != nil {
		return id, err
	}

	// Apply the fixes
	if err := c.fixContainer(id); err != nil {
		return id, err
	}
	check(err)

	// Clean up on ctrl+c
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	signal.Notify(ch, syscall.SIGTERM)

	go func() {
		<-ch
		log.Debug("Trapped ctrl+c")

		if err = c.DeleteContainer(id); err != nil {
			log.Errorf("Failed to remove container: %s", err)
		}
		os.Exit(1)
	}()

	// Execute the command
	err = c.execContainer(id, c.Entrypoint, "user", false)
	if err != nil {
		return id, err
	}

	// Container has finished running. Get its exit code
	inspect, err := c.Cli.ContainerInspect(context.Background(), id)
	if err != nil {
		return id, err
	}

	// Delete the container if required
	if rm {
		if err = c.DeleteContainer(id); err != nil {
			return id, err
		}
	}

	if inspect.State.ExitCode != 0 {
		return id, err
	}
	return id, nil
}

// ContainerExists determines if a container with the passed id exists
func (c *DockerClient) ContainerExists(id string) bool {
	_, err := c.Cli.ContainerInspect(context.Background(), id)

	return err == nil
}

// DeleteContainer does exactly that adding some logging
func (c *DockerClient) DeleteContainer(id string) error {
	log.WithFields(log.Fields{
		"id": id[0:12],
	}).Debug("Removing container")

	if err := c.Cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{Force: true}); err != nil {
		return err
	}
	return nil
}

// fixContainer executes as script inside a running container that sets it up
// for the user:
// - It adds the current user in the container.
// - It adds the current user main group and sets up the user on that one.
// - It creates user home (used by certain tools to store cache files)
func (c *DockerClient) fixContainer(containerID string) error {

	log.WithField("image", c.Conf.Image).Debug("Fixing image")

	// To deploy the fix script we need to tar it first
	tar := new(archivex.TarFile)
	tar.Create("/tmp/clide_fix.tar")
	dat, err := ioutil.ReadFile("../static/fix.sh")
	check(err)
	tar.Add("fix.sh", dat)
	tar.Close()
	holyTar, err := os.Open("/tmp/clide_fix.tar")
	defer holyTar.Close()

	// The actual deploy happens next
	err = c.Cli.CopyToContainer(context.Background(),
		containerID,
		"/tmp",
		holyTar,
		types.CopyToContainerOptions{})
	if err != nil {
		return err
	}

	// Executing the fix script.
	err = c.execContainer(
		containerID,
		[]string{"/bin/sh", "/tmp/fix.sh"},
		"root",
		true)
	if err != nil {
		return err
	}

	return err
}

// initializeContainer run the container with a background void process
// such as 'wait infinite'
func (c *DockerClient) initializeContainer(name string) (string, error) {

	log.WithFields(log.Fields{
		"image": c.Conf.Image,
		"envs":  fmt.Sprintf("%v", c.Conf.Env),
		"cmd":   fmt.Sprintf("%v", c.Conf.Cmd),
	}).Debug("Initializing new container")

	// Pulling the Image
	if err := c.PullImage(c.Conf.Image); err != nil {
		return "", fmt.Errorf("Failed to fetch image: %s", err)
	}

	// Creation of the container
	resp, err := c.Cli.ContainerCreate(context.Background(), c.Conf, c.HostConf,
		c.NetConf, name)
	if err != nil {
		return "", fmt.Errorf("Failed to create container: %s", err)
	}

	// Starting the container
	if err := c.Cli.ContainerStart(context.Background(), resp.ID,
		types.ContainerStartOptions{}); err != nil {
		return resp.ID, err
	}

	out := int(os.Stdout.Fd())
	// Resizing the container output
	tw, th, _ := terminal.GetSize(out)
	err = c.Cli.ContainerResize(context.Background(), resp.ID,
		types.ResizeOptions{Height: uint(th), Width: uint(tw)})
	if err != nil {
		return resp.ID, err
	}

	return resp.ID, nil
}

// execContainer executes a command inside a running container
func (c *DockerClient) execContainer(id string, cmd []string,
	user string, nonInteractive bool) error {

	log.WithFields(log.Fields{
		"id":   id[0:12],
		"cmd":  fmt.Sprintf("%v", cmd),
		"user": user,
	}).Debug("Executing command in container as user")

	// Creating execution for the entrypoint shell script.
	resp, err := c.Cli.ContainerExecCreate(context.Background(), id,
		types.ExecConfig{
			Cmd:          cmd,
			Tty:          true,
			AttachStdout: true,
			AttachStderr: true,
			AttachStdin:  true,
			User:         user,
			Detach:       false,
		})
	if err != nil {
		panic(err)
	}

	execID := resp.ID

	in := int(os.Stdin.Fd())

	// Attaching to the exec
	hijack, err := c.Cli.ContainerExecAttach(context.Background(), execID,
		types.ExecStartCheck{Tty: true})
	if err == nil {
		defer hijack.Conn.Close()
	}
	if err != nil {
		panic(err)
	}

	if !nonInteractive && terminal.IsTerminal(int(os.Stdin.Fd())) {

		log.Debug("Running interactively")
		// While we have a container running, create a buffer for the pscli logs
		logBuffer := bufio.NewWriter(os.Stdout)
		log.SetOutput(logBuffer)
		// Write buffer to stdout once detatched from container
		defer logBuffer.Flush()
		// Reset logs to stdout after conection is closed
		defer log.SetOutput(os.Stdout)

		// Making the terminal raw
		oldState, err := terminal.MakeRaw(in)
		if err != nil {
			panic(err)
		}
		defer terminal.Restore(in, oldState)

		// Start stdin reader
		go func() {
			log.Debug("Listening to stdin")

			if _, err := io.Copy(hijack.Conn, os.Stdin); err != nil {
				log.Errorf("Write error: %s", err)
			}
		}()
	}

	// Start stdout writer
	if _, err := io.Copy(os.Stdout, hijack.Conn); err != nil {
		log.Errorf("Read error: %s", err)
	}

	return err
}
