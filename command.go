package cali

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Command is the actual command run by the cli and essentially just wraps cobra.Command and
// has an associated Task
type Command struct {
	name    string
	RunTask *Task
	cobra   *cobra.Command
}

// NewCommand returns an freshly initialised command
func NewCommand(n string) *Command {
	return &Command{
		name:  n,
		cobra: &cobra.Command{Use: n},
	}
}

// SetShort sets the short description of the command
func (c *Command) SetShort(s string) {
	c.cobra.Short = s
}

// SetLong sets the long description of the command
func (c *Command) SetLong(l string) {
	c.cobra.Long = l
}

// setPreRun sets the cobra.Command.PreRun function
func (c *Command) setPreRun(f cobraFunc) {
	c.cobra.PreRun = f
}

// setRun sets the cobra.Command.Run function
func (c *Command) setRun(f cobraFunc) {
	c.cobra.Run = f
}

// AddTask is something executed by a command
func (c *Command) AddTask(def interface{}) *Task {
	t := NewTask()

	switch d := def.(type) {
	case string:
		t.SetImage(d)
	case TaskFunc:
		t.SetFunc(d)
	default:
		// Slightly unidiomatic to blow up here rather than return an error
		// choosing to so as to keep the API uncluttered and also if you get here it's
		// an implementation error rather than a runtime error.
		fmt.Println("Unknown Task type. Must either be an image (string) or a TaskFunc")
		os.Exit(exitCodeAPIError)
	}
	c.RunTask = t
	return t
}

// Flags returns the FlagSet for the command and is used to set new flags for the command
func (c *Command) Flags() *pflag.FlagSet {
	return c.cobra.PersistentFlags()
}

// BindFlags needs to be called after all flags for a command have been defined
func (c *Command) BindFlags() {
	c.Flags().VisitAll(func(f *pflag.Flag) {
		flags.BindPFlag(f.Name, f)
		flags.SetDefault(f.Name, f.DefValue)
	})
}
