# cali
--
    import "github.com/edupo/cali"


## Usage

#### type Cli

```go
type Cli struct {
	*Command
}
```

Cli is the application itself

#### func  NewCli

```go
func NewCli(n string) *Cli
```
NewCli returns a brand new cli

#### func (*Cli) AddCommand

```go
func (c *Cli) AddCommand(n string) *Command
```
AddCommand returns a brand new command attached to it's parent cli

#### func (*Cli) FlagValues

```go
func (c *Cli) FlagValues() *viper.Viper
```
FlagValues returns the wrapped viper object allowing the API consumer to use
methods like GetString to get values from config

#### func (*Cli) Start

```go
func (c *Cli) Start()
```
Start the fans please!

#### type Command

```go
type Command struct {
	RunTask *Task
}
```

Command is the actual command run by the cli and essentially just wraps
cobra.Command and has an associated Task

#### func  NewCommand

```go
func NewCommand(n string) *Command
```
NewCommand returns an freshly initialised command

#### func (*Command) AddTask

```go
func (c *Command) AddTask(def interface{}) *Task
```
AddTask is something executed by a command

#### func (*Command) BindFlags

```go
func (c *Command) BindFlags()
```
BindFlags needs to be called after all flags for a command have been defined

#### func (*Command) Flags

```go
func (c *Command) Flags() *pflag.FlagSet
```
Flags returns the FlagSet for the command and is used to set new flags for the
command

#### func (*Command) SetLong

```go
func (c *Command) SetLong(l string)
```
SetLong sets the long description of the command

#### func (*Command) SetShort

```go
func (c *Command) SetShort(s string)
```
SetShort sets the short description of the command

#### type Task

```go
type Task struct {
	*docker.Client
}
```

Task is the action performed when it's parent command is run

#### func  NewTask

```go
func NewTask() *Task
```
NewTask returns a new Task structure containing a new Client object.

#### func (*Task) Bind

```go
func (t *Task) Bind(src, dst string) (string, error)
```
Bind is a utility function which will return the correctly formatted string when
given a source and destination directory

The ~ symbol and relative paths will be correctly expanded depending on the host
OS

#### func (*Task) BindDocker

```go
func (t *Task) BindDocker()
```
BindDocker - Task util (convenience) to Bind the docker socket.

#### func (*Task) SetDefaults

```go
func (t *Task) SetDefaults(args []string) error
```
SetDefaults sets the default host config for a task container Mounts the PWD to
/tmp/workspace Mounts your ~/.aws directory to /root - change this if your image
runs as a non-root user Sets /tmp/workspace as the workdir Configures git

#### func (*Task) SetFunc

```go
func (t *Task) SetFunc(f TaskFunc)
```
SetFunc sets the TaskFunc which is run when the parent command is run if this is
left unset, the defaultTaskFunc will be executed instead

#### func (*Task) SetInitFunc

```go
func (t *Task) SetInitFunc(f TaskFunc)
```
SetInitFunc sets the TaskFunc which is executed before the main TaskFunc. It's
pupose is to do any setup of the Client which depends on command line args for
example

#### type TaskFunc

```go
type TaskFunc func(t *Task, args []string)
```

TaskFunc is a function executed by a Task when the command the Task belongs to
is run
