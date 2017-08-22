package cli

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"geeny/output"
	"geeny/util"
	"geeny/version"
)

// Output encapsulates the command line output
type Output struct {
}

// NewOutput returns a new Output instance
func NewOutput() *Output {
	return &Output{}
}

const (
	cmds = " [subcommands]"
	opts = " [options]"
)

// - private

func (o *Output) printCompletionForCommand(cmd *Command) *errorContext {
	if len(cmd.Commands) == 0 && len(cmd.Options) == 0 {
		return nil
	}

	names := []string{}
	for _, c := range cmd.Commands {
		names = append(names, c.Name)
	}
	for _, o := range cmd.Options {
		names = append(names, "-"+o.Flag)
		for _, a := range o.Aliases {
			names = append(names, "--"+a)
		}
	}

	completion := names[0]
	for _, n := range names[1:] {
		completion += " " + n
	}
	output.Println(completion)

	return &errorContext{
		internalErr: errors.New(""),
		errCode:     errBashCompletion,
	}
}

func (o *Output) printUsage(cmd *Command) *errorContext {
	if cmd.Hidden {
		return &errorContext{
			internalErr: errors.New("command not supported"),
			errCode:     errCommandNotSupported,
		}
	}

	output.Println(version.Version)
	o.printCommandInfo(cmd, true, 0)

	if len(cmd.Options) > 0 {
		output.Println("\n\x1b[39;1moptions:\x1b[0m")
		o.printOptionInfos(cmd.Options)
	}

	if len(cmd.Commands) > 0 {
		output.Println("\n\x1b[39;1msubcommands:\x1b[0m")
		o.printCommandInfos(cmd.Commands)
	}

	return &errorContext{
		internalErr: errors.New(""),
		errCode:     errUsage,
	}
}

func (o *Output) printCommandInfo(cmd *Command, top bool, minlen int) {
	if cmd.Hidden {
		return
	}

	var preText string
	var splitText string
	var finalText string
	var comment string

	if top {
		preText = o.getPreText(cmd)
		splitText = "\n"
		if len(cmd.Comment) == 0 {
			if len(cmd.Summary) > 0 {
				splitText = splitText + "\n"
				finalText = "\n"
			}
			comment = ""
		} else {
			finalText = ".\n"
			comment = finalText
		}
	} else {
		preText = ""
		if len(cmd.Summary) == 0 {
			splitText = " \x1b[39;1m "
		} else {
			splitText = " \x1b[39;1m-\x1b[0m "
		}
		finalText = "\n"
		if len(cmd.Comment) == 0 {
			comment = ""
		} else {
			comment = "\n" + strings.Repeat(" ", minlen+3)
		}
	}
	output.Printf("%s%s%s%s%s%s", preText, o.getCommandInfo(cmd, minlen, true), splitText, cmd.Summary, comment+cmd.Comment, finalText)
}

func (o *Output) printCommandInfos(cmd []*Command) {
	var namelen = 0
	for _, c := range cmd {
		namelen = util.Max(namelen, len(o.getCommandInfo(c, 0, false))+len(o.getCommandOptionsInfo(c)))
	}
	for _, c := range cmd {
		o.printCommandInfo(c, false, namelen)
	}
}

func (o *Output) printOptionInfos(options []*Option) {
	var namelen = 0
	for _, option := range options {
		namelen = util.Max(namelen, len(o.getOptionInfo(option, 0)))
	}
	for _, option := range options {
		output.Printf("%s\x1b[39;1m-\x1b[0m %s\n", o.getOptionInfo(option, namelen), option.Description)
	}
}

func (o *Output) printOption(option *Option) {
	output.Printf("\x1b[39;1m%s:\x1b[0m ", option.Description)
}

func (o *Output) getPreText(cmd *Command) string {
	var stack []*Command
	for p := cmd.Parent; p != nil; p = p.Parent {
		stack = append(stack, p)
	}
	var buf bytes.Buffer
	buf.WriteString("\x1b[39;1mUsage:\x1b[0m ")
	for i := len(stack) - 1; i >= 0; i-- {
		buf.WriteString(stack[i].Name)
		buf.WriteString(" ")
	}
	return buf.String()
}

func (o *Output) getCommandInfo(cmd *Command, minlen int, withOptions bool) string {
	if withOptions {
		var opts = o.getCommandOptionsInfo(cmd)
		var pad = minlen - len(cmd.Name) - len(opts)
		var padding = ""
		if pad > 0 {
			padding = strings.Repeat(" ", pad)
		}
		return fmt.Sprintf("%s\x1b[39;1m%s\x1b[0m%s", cmd.Name, opts, padding)
	}
	if minlen > 0 {
		return fmt.Sprintf("%-"+strconv.Itoa(minlen), cmd.Name)
	}
	return cmd.Name
}

func (o *Output) getCommandOptionsInfo(cmd *Command) string {
	var cmdsPresent = len(cmd.Commands) > 0
	var optsPresent = len(cmd.Options) > 0
	if cmdsPresent && optsPresent {
		return cmds + opts
	}
	if cmdsPresent {
		return cmds
	}
	if optsPresent {
		return opts
	}
	return ""
}

func (o *Output) getOptionInfo(option *Option, minlen int) string {
	var name string
	if option.Type == OptionTypeBool {
		name = ""
	} else {
		name = " <" + option.Name + ">"
	}

	aliases := ""
	for _, alias := range option.Aliases {
		aliases = aliases + ", --" + alias
	}

	var namelen = util.Max(minlen, len(name))
	return fmt.Sprintf("-%-"+strconv.Itoa(namelen)+"s", option.Flag+aliases+name)
}
