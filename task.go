package cali

import (
	"os/user"
	"strings"
	"path/filepath"
	"fmt"
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
//  - Sets /tmp/workspace as the workdir
//  - Publishes HOST_USER_ID and HOST_GROUP_ID in the container
//  - Mounts the PWD to /tmp/workspace
//  - Configures git
//  - Set the default command
func (t *Task) SetDefaults(args []string) error {
	t.SetWorkDir(workDir)

	u, err := user.Current()
	if err != nil {
		return err
	}
	t.AddEnv("HOST_USER_ID", u.Uid)
	t.AddEnv("HOST_GROUP_ID", u.Gid)

	err = t.BindFromGit(gitCfg, func() error {
		err := t.Bind("./", workDir)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	t.SetCmd(args)
	return nil
}

// BindDocker - Task util (convenience) to Bind the docker socket.
func (t *Task) BindDocker() error {
	err := t.Bind("/var/run/docker.sock", "/var/run/docker.sock")
	if err != nil {
		return err
	}
	return nil
}

// Bind - Task util to add a Bind. '~' in src will be expanded according to the current user for convenience.
func (t *Task) Bind(src, dst string) error {
	var expanded string

	if strings.HasPrefix(src, "~") {
		usr, err := user.Current()
		if err != nil {
			return err
		}
		expanded = filepath.Join(usr.HomeDir, src[2:])
	} else {
		expanded = src
	}

	expanded, err := filepath.Abs(expanded)
	if err != nil {
		return err
	}

	t.AddBind(fmt.Sprintf("%s:%s", expanded, dst))
	return nil
}