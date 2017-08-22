package cli

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	log "geeny/log"
	"reflect"
)

var commandNotFoundError = "could not find command for path %v. please update the cli, and if this error still shows, report a bug to support@geeny.io"

// Command represents a cli command
type Command struct {
	Name        string
	Summary     string
	Comment     string
	Help        string
	Interactive bool
	Action      func(context *Context) (*Meta, error)
	Options     []*Option
	Commands    []*Command
	Parent      *Command
	Hidden      bool
	Extension   interface{}
	NonCategory bool
}

// CommandForPath returns command for given path or nil if it doesnt exist
func (c *Command) CommandForPath(path []string, idx int) (*Command, error) {
	if idx >= len(path) {
		return nil, fmt.Errorf(commandNotFoundError, path)
	}
	// check command
	p := path[idx]
	if strings.Compare(c.Name, p) != 0 {
		return nil, fmt.Errorf(commandNotFoundError, path)
	}
	if idx == len(path)-1 {
		return c, nil
	}
	// check subcommands
	for _, sub := range c.Commands {
		result, _ := sub.CommandForPath(path, idx+1)
		if result != nil {
			return result, nil
		}
	}
	return nil, fmt.Errorf(commandNotFoundError, path)
}

// SubCommandForName returns sub command for name or nil if it doesnt exist
func (c *Command) SubCommandForName(name string) *Command {
	for _, command := range c.Commands {
		if command.Name == name {
			return command
		}
	}
	return nil
}

// SubCommandForPath returns sub command for underscore separated path, e.g. my_path_is_this
func (c *Command) SubCommandForPath(path []string) (*Command, bool) {
	cmd := c
	for i := 0; i < len(path); i++ {
		cmd = cmd.SubCommandForName(path[i])
		if cmd == nil {
			return nil, false
		}
	}
	return cmd, true
}

// OptionForFlag returns an options for provided flag
func (c *Command) OptionForFlag(flag string) (*Option, error) {
	for _, o := range c.Options {
		if strings.Compare(o.Flag, flag) == 0 {
			return o, nil
		}
	}
	return nil, errors.New("can't find option flag: " + flag + ", for command: " + c.Name)
}

// SetValueForOptionWithFlag sets value for an option that has given flag
func (c *Command) SetValueForOptionWithFlag(value interface{}, flag string) error {
	if reflect.TypeOf(value).Kind() != reflect.Ptr {
		return errors.New("command " + c.Name + ", " + flag + ".Value should be a pointer")
	}
	o, err := c.OptionForFlag(flag)
	if err != nil {
		return err
	}
	o.Value = value
	return nil
}

// Exec executes a given command
func (c *Command) Exec() (meta *Meta, err error) {
	ctx := new(Context)
	ctx.Command = c
	for _, o := range c.Options {
		ctx.Args = append(ctx.Args, o)
	}
	return c.Action(ctx)
}

// Merge sub commands into a root command
func (c *Command) Merge(subCommand *Command) error {
	command := c.SubCommandForName(subCommand.Name)
	if command != nil {
		if len(subCommand.Commands) == 0 {
			if command.NonCategory || !subCommand.NonCategory {
				log.Infof("Ignore %v -> %v", c.Name, subCommand.Name)
				return nil
			}
			log.Infof("Overwrite %v -> %v", c.Name, subCommand.Name)
			command.NonCategory = subCommand.NonCategory
			command.Extension = subCommand.Extension
			command.Hidden = subCommand.Hidden
			command.Summary = subCommand.Summary
			command.Comment = subCommand.Comment
			command.Help = subCommand.Help
			command.Action = subCommand.Action
			command.Interactive = subCommand.Interactive
			command.Options = subCommand.Options
			return nil
		}
		log.Tracef("Merge %v -> %v", c.Name, subCommand.String())
		for _, cmd := range subCommand.Commands {
			err := command.Merge(cmd)
			if err != nil {
				return err
			}
		}
		return nil
	}
	log.Tracef("Append %v -> %v", c.Name, subCommand.String())
	c.Commands = append(c.Commands, subCommand)
	return nil
}

type CommandsAlpha []*Command

func (a CommandsAlpha) Len() int           { return len(a) }
func (a CommandsAlpha) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CommandsAlpha) Less(i, j int) bool { return a[i].Name < a[j].Name }

func (c *Command) Sort() {
	sort.Stable(OptionsAlpha(c.Options))
	sort.Stable(CommandsAlpha(c.Commands))
	for _, s := range c.Commands {
		s.Sort()
	}
}

func (cmd *Command) SetupParents() {
	var hidden = true
	for i := range cmd.Commands {
		hidden = hidden && cmd.Commands[i].Hidden
		cmd.Commands[i].Parent = cmd
		cmd.Commands[i].SetupParents()
	}

	// This is to hide non-leaf commands that only have hidden subcommands and no options
	if hidden && len(cmd.Commands) > 0 && len(cmd.Options) <= 0 {
		cmd.Hidden = true
	}
}

func (c *Command) String() string {
	var sub []string
	for _, n := range c.Commands {
		sub = append(sub, n.Name)
	}
	return fmt.Sprintf("%v (%v): %v", c.Name, strings.Join(sub, " "), c.Extension)
}

// TreeString returns command as a string tree, e.g.
/*
1a
  2a
    3a
  2b
1b
1c
  2
*/
func (c *Command) TreeString(level int) string {
	str := fmt.Sprintf("%s - %s", c.Name, c.Summary)
	for _, sub := range c.Commands {
		str = fmt.Sprintf("%s\n%s%s", str, strings.Repeat(" ", level+1), sub.TreeString(level+1))
	}
	return str
}
