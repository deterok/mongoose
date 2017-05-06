package mongoose

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	flag "github.com/spf13/pflag"
)

type Command struct {
	Name string

	MinArgs int
	Args    []string
	Tail    []string

	PPreRun  func(c *Command)
	PreRun   func(c *Command)
	Run      func(c *Command)
	PostRun  func(c *Command)
	PPostRun func(c *Command)

	commands map[string]*Command
	flags    *flag.FlagSet
	output   io.Writer
}

func (c *Command) Flags() *flag.FlagSet {
	if c.flags == nil {
		c.flags = flag.NewFlagSet(c.Name, flag.ContinueOnError)
	}

	return c.flags
}

func (c *Command) AddCommand(command *Command) {
	c.AddNamedCommand(strings.ToLower(command.Name), command)
}

func (c *Command) AddNamedCommand(name string, command *Command) {
	if c.commands == nil {
		c.commands = make(map[string]*Command)
	}

	c.commands[name] = command
	c.Flags().SetInterspersed(false)

	if c.Output() != nil && command.Output() == nil {
		command.SetOutput(c.Output())
	}
}

func (c *Command) Execute(args []string) {
	err := c.Parse(args)
	if err != nil {
		return
	}

	if c.PPreRun != nil {
		c.PPreRun(c)
	}

	if len(c.Tail) > 0 {
		subCMD := c.GetCommand(c.Tail[0])
		if subCMD != nil {
			subCMD.Execute(c.Tail[1:])
		}

	} else {
		c.ExecuteCmd()
	}

	if c.PPostRun != nil {
		c.PPostRun(c)
	}
}

func (c *Command) Parse(args []string) error {
	if err := c.Flags().Parse(args); err != nil {
		return err
	}
	parsedArgs := c.Flags().Args()

	switch {
	case c.MinArgs > len(parsedArgs):
		errString := fmt.Sprintf("The number of incoming arguments is less than the minimum (%d' > %d)",
			c.MinArgs, len(parsedArgs))
		return errors.New(errString)

	case c.MinArgs == 0:
		c.Tail = parsedArgs

	case c.MinArgs > 0:
		c.Args = parsedArgs[:c.MinArgs]
		c.Tail = parsedArgs[c.MinArgs:]

	case c.MinArgs < 0:
		c.Args = parsedArgs
	}

	return nil
}

func (c *Command) ExecuteCmd() {
	if c.PreRun != nil {
		c.PreRun(c)
	}

	if c.Run != nil {
		c.Run(c)
	}

	if c.PostRun != nil {
		c.PostRun(c)
	}
}

func (c *Command) GetCommand(name string) *Command {
	cmd := c.commands[name]
	return cmd
}

func (c *Command) FindCommandByPath(path, sep string) *Command {
	names := strings.Split(path, sep)
	return c.FindCommand(names)
}

func (c *Command) FindCommand(names []string) *Command {
	if len(names) == 0 {
		return c
	}

	cmd := c.GetCommand(names[0])
	if cmd != nil {
		return cmd.FindCommand(names[1:])
	}

	return cmd
}

func (c *Command) SetOutput(output io.Writer) {
	c.output = output
	c.Flags().SetOutput(output)
}

func (c *Command) Output() io.Writer {
	if c.output == nil {
		return os.Stderr
	}

	return c.output
}
