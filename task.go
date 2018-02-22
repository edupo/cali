package cali

import (
	"fmt"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

// Task is the action performed when it's parent command is run
type Task struct {
	f, init TaskFunc
	*DockerClient
}

// SetFunc sets the TaskFunc which is run when the parent command is run
// if this is left unset, the defaultTaskFunc will be executed instead
func (t *Task) SetFunc(f TaskFunc) {
	t.f = f
}

// SetInitFunc sets the TaskFunc which is executed before the main TaskFunc. It's
// pupose is to do any setup of the DockerClient which depends on command line args
// for example
func (t *Task) SetInitFunc(f TaskFunc) {
	t.init = f
}

// SetDefaults sets the default host config for a task container
// Mounts the PWD to /tmp/workspace
// Mounts your ~/.aws directory to /root - change this if your image runs as a non-root user
// Sets /tmp/workspace as the workdir
// Configures git
func (t *Task) SetDefaults(args []string) error {
	t.SetWorkDir(workDir)
	awsDir, err := t.Bind("~/.aws", "/root/.aws")
	if err != nil {
		return err
	}
	t.AddBinds([]string{awsDir})

	err = t.BindFromGit(gitCfg, func() error {
		pwd, err := t.Bind("./", workDir)
		if err != nil {
			return err
		}
		t.AddBinds([]string{pwd})
		return nil
	})
	if err != nil {
		return err
	}
	t.SetCmd(args)
	return nil
}

// Bind is a utility function which will return the correctly formatted string when given a source
// and destination directory
//
// The ~ symbol and relative paths will be correctly expanded depending on the host OS
func (t *Task) Bind(src, dst string) (string, error) {
	var expanded string

	if strings.HasPrefix(src, "~") {
		usr, err := user.Current()

		if err != nil {
			return expanded, fmt.Errorf("Error expanding bind path: %s", src)
		}
		expanded = filepath.Join(usr.HomeDir, src[2:])
	} else {
		expanded = src
	}
	expanded, err := filepath.Abs(expanded)

	if err != nil {
		return expanded, fmt.Errorf("Error expanding bind path: %s", src)
	}
	return fmt.Sprintf("%s:%s", expanded, dst), nil
}

// BindDocker - Task util (convenience) to Bind the docker socket.
func (t *Task) BindDocker() {
	dockerSocket := "/var/run/docker.sock"
	if runtime.GOOS == "windows" {
		dockerSocket = "//var/run/docker.sock"
	}
	str, err := t.Bind(dockerSocket, "/var/run/docker.sock")
	if err != nil {
		panic(err)
	}
	t.AddBind(str)
}
