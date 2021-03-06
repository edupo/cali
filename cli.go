package cali

import (
	"fmt"
	"os"
	"runtime"

	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/edupo/cali/docker"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	exitCodeRuntimeError = 1
	exitCodeAPIError     = 2

	defaultDockerRegistry = "docker.io"
	workDir               = "/tmp/workspace"
)

var (
	debug, jsonLogs, nonInteractive bool
	dockerHost, dockerRegistry      string
	flags                           *viper.Viper
	gitCfg                          *docker.GitCheckoutConfig
	cwd                             string
)

// cobraFunc represents the function signiture which cobra uses for it's Run, PreRun, PostRun etc.
type cobraFunc func(cmd *cobra.Command, args []string)

// commands is a set of commands
type commands map[string]*Command

// Cli is the application itself
type Cli struct {
	name    string
	cfgFile *string
	cmds    commands
	*Command
}

// NewCli returns a brand new cli
func NewCli(n string) *Cli {
	c := Cli{
		name:    n,
		cmds:    make(commands),
		Command: NewCommand(n),
	}
	c.cobra.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if debug {
			log.SetLevel(log.DebugLevel)
		}

		if jsonLogs {
			log.SetFormatter(&log.JSONFormatter{})
		}
	}
	flags = viper.New()
	return &c
}

// AddCommand returns a brand new command attached to it's parent cli
func (c *Cli) AddCommand(n string) *Command {
	cmd := NewCommand(n)
	c.cmds[n] = cmd
	cmd.setPreRun(func(c *cobra.Command, args []string) {
		if cmd.RunTask.init != nil {
			cmd.RunTask.init(cmd.RunTask, args)
		}
	})
	cmd.setRun(func(c *cobra.Command, args []string) {
		cmd.RunTask.f(cmd.RunTask, args)
	})
	c.cobra.AddCommand(cmd.cobra)
	return cmd
}

// FlagValues returns the wrapped viper object allowing the API consumer to use methods
// like GetString to get values from config
func (c *Cli) FlagValues() *viper.Viper {
	return flags
}

// initFlags does the intial setup of the root command's persistent flags
func (c *Cli) initFlags() {
	var cfg string
	txt := fmt.Sprintf("config file (default is $HOME/.%s.yaml)", c.name)
	c.cobra.PersistentFlags().StringVar(&cfg, "config", "", txt)
	c.cfgFile = &cfg

	wd, _ := os.Getwd()
	c.Flags().StringVarP(&cwd, "work-dir", "C", wd,
		"URI of Docker Daemon")
	flags.BindPFlag("work-dir", c.Flags().Lookup("work-dir"))

	var dockerSocket string
	if runtime.GOOS == "windows" {
		dockerSocket = "npipe:////./pipe/docker_engine"
	} else {
		dockerSocket = "unix:///var/run/docker.sock"
	}
	c.Flags().StringVarP(&dockerHost, "docker-host", "H", dockerSocket, "URI of Docker Daemon")
	flags.BindPFlag("docker-host", c.Flags().Lookup("docker-host"))
	flags.SetDefault("docker-host", dockerSocket)

	c.Flags().StringVarP(&dockerRegistry, "docker-registry", "R", defaultDockerRegistry, "URI of Docker registry")
	flags.BindPFlag("docker-registry", c.Flags().Lookup("docker-registry"))
	flags.SetDefault("docker-registry", defaultDockerRegistry)

	c.Flags().BoolVarP(&debug, "debug", "d", false, "Debug mode")
	flags.BindPFlag("debug", c.Flags().Lookup("debug"))
	flags.SetDefault("debug", true)

	c.Flags().BoolVarP(&jsonLogs, "json", "j", false, "Log in json format")
	flags.BindPFlag("json", c.Flags().Lookup("json"))
	flags.SetDefault("json", true)

	c.Flags().BoolVarP(&nonInteractive, "non-interactive", "N", false, "Do not create a tty for Docker")
	flags.BindPFlag("non-interactive", c.Flags().Lookup("non-interactive"))
	flags.SetDefault("non-interactive", false)

	gitCfg = new(docker.GitCheckoutConfig)
	c.Flags().StringVarP(&gitCfg.Repo, "git", "g", "", "Git repo to checkout and build. Default behaviour is to build $PWD.")
	flags.BindPFlag("git", c.Flags().Lookup("git"))

	c.Flags().StringVarP(&gitCfg.Branch, "git-branch", "b", "master", "Branch to checkout. Only makes sense when combined with the --git flag.")
	flags.BindPFlag("branch", c.Flags().Lookup("branch"))
	flags.SetDefault("branch", "master")

	c.Flags().StringVarP(&gitCfg.RelPath, "git-path", "P", "", "Path within a git repo where we want to operate.")
	flags.BindPFlag("git-path", c.Flags().Lookup("git-path"))
}

// initConfig does the initial setup of viper
func (c *Cli) initConfig() {
	if *c.cfgFile != "" {
		flags.SetConfigFile(*c.cfgFile)
	} else {
		flags.SetConfigName(fmt.Sprintf(".%s", c.name))
		flags.AddConfigPath("$HOME")
		flags.AddConfigPath(".")
	}
	flags.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	flags.AutomaticEnv()

	// If a config file is found, read it in
	if err := flags.ReadInConfig(); err == nil {
		log.WithField("file", flags.ConfigFileUsed()).Info(
			"Using configuration file")
	}
}

// Start the fans please!
func (c *Cli) Start() {
	c.initFlags()
	cobra.OnInitialize(c.initConfig)

	if err := c.cobra.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(exitCodeRuntimeError)
	}
}
