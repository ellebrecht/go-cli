package cli

import (
	"bytes"
	"errors"
	"flag"
	"fmt"

	output "geeny/output"
)

// - private

type parser struct {
	output           *Output
	isActionsEnabled bool
}

func newParser(isActionsEnabled bool) *parser {
	return &parser{
		output:           NewOutput(),
		isActionsEnabled: isActionsEnabled,
	}
}

func (p *parser) parseCommand(cmd *Command, args []string, index int) (*Meta, *errorContext) {
	// if there's nothing else after the current argument, attempt to print usage / do something interactive
	if index >= len(args) {
		return p.performUsageAction(cmd)
	}

	// find the command to parse for the given argument
	arg := args[index]
	argCmd := cmd
	if index > 0 {
		argCmd = p.commandForArg(cmd.Commands, arg)
	}

	// if we don't find a related command
	if argCmd == nil {
		// make sure it's not a flag like -h or --help
		_, err := p.parseOptions(cmd, args[index:])
		if err != nil {
			return nil, p.performErrorAction(cmd, err)
		}

		// else it's not supported
		return nil, &errorContext{
			internalErr: errors.New("command not supported: " + arg),
			errCode:     errCommandNotSupported,
		}
	}

	// if the command has both possible subcommands and options...
	if len(argCmd.Commands) > 0 && len(argCmd.Options) > 0 {
		// we first try parsing the next command
		meta, cmdErr := p.parseCommand(argCmd, args, index+1)
		if cmdErr == nil { // if we succeeded, return without an error
			return meta, nil
		}
		// if the command is not supported, it might be an option (so we continue), else we report the error
		if cmdErr.errCode != errCommandNotSupported {
			return nil, cmdErr
		}
		// otherwise, return the result of parsing the options and executing the action
		return p.parseOptionsAndPerformAction(argCmd, args[index+1:])
	}

	// if we only expect more commands, and there are arguments left on the command line OR there is no action, try to parse the next subcommand
	if len(argCmd.Commands) > 0 && (index+1 < len(args) || argCmd.Action == nil) {
		return p.parseCommand(argCmd, args, index+1)
	}

	// if we only expect options, but none are provided, attempt to do something interactive (like getting the flags) or just print usage
	if p.numOfMissingOptions(argCmd.Options) > 0 && index+1 >= len(args) {
		return p.performUsageAction(argCmd)
	}

	// if we only expect options, lets extract the -flag information and execute the command's action with this as context
	// or if there's extra stuff on the command line, see if it's important
	if len(argCmd.Options) > 0 || index < len(args) {
		return p.parseOptionsAndPerformAction(argCmd, args[index+1:])
	}

	// if we reached here, the command doesn't require context (i.e. -flags), so just execute its action
	return p.performAction(argCmd, nil)
}

func (p *parser) parseOptionsAndPerformAction(cmd *Command, args []string) (*Meta, *errorContext) {
	ctx, err := p.parseOptions(cmd, args)
	if err != nil {
		return nil, p.performErrorAction(cmd, err)
	}
	return p.performAction(cmd, ctx)
}

func (p *parser) parseOptions(cmd *Command, args []string) (*Context, *errorContext) {
	context := new(Context)
	context.Command = cmd
	flagSet := flag.NewFlagSet(cmd.Name, flag.ContinueOnError)

	// override default usage message & err output as we handle it
	flagSet.Usage = func() {}
	flagSet.SetOutput(bytes.NewBuffer([]byte{}))

	// setup parsing (flags, aliases, etc)
	for _, option := range cmd.Options {
		option.possibleValues = []interface{}{}     // reset any previous attempts (as we're using pointers)
		allFlags := option.allFlags()               // flag + aliases
		context.Args = append(context.Args, option) // add to context

		switch option.Type {
		case OptionTypeBool:
			for _, alias := range allFlags {
				defaultVal := option.defaultValue().(bool)
				if option.DefaultValue != nil {
					option.Value = &defaultVal
				}
				option.possibleValues = append(option.possibleValues, flagSet.Bool(alias, defaultVal, option.Description))
			}
		case OptionTypeInt:
			for _, alias := range allFlags {
				defaultVal := option.defaultValue().(int)
				if option.DefaultValue != nil {
					option.Value = &defaultVal
				}
				option.possibleValues = append(option.possibleValues, flagSet.Int(alias, defaultVal, option.Description))
			}
		case OptionTypeString:
			for _, alias := range allFlags {
				defaultVal := option.defaultValue().(string)
				if option.DefaultValue != nil {
					option.Value = &defaultVal
				}
				option.possibleValues = append(option.possibleValues, flagSet.String(alias, defaultVal, option.Description))
			}
		default:
			panic("unhandled switch statement")
		}
	}

	// add bash completion option
	bashCompletionOption := Option{
		Flag: FlagGenerateBashCompletion,
	}
	bashCompletionOption.Value = flagSet.Bool(bashCompletionOption.Flag, false, "")
	//context.Args = append(context.Args, &bashCompletionOption) //TODO why?

	// parse
	err := flagSet.Parse(args)
	if err != nil {
		return nil, &errorContext{
			internalErr: err,
			errCode:     errParseOptions,
		}
	}

	// if bash completion flag detected, return error
	if *bashCompletionOption.Value.(*bool) {
		return nil, &errorContext{
			errCode: errBashCompletion,
		}
	}

	// check for unknown commands
	if flagSet.NArg() > 0 {
		return nil, &errorContext{
			internalErr: fmt.Errorf("found an unknown option somewhere around where you see: %s", flagSet.Args()),
			errCode:     errParseOptions,
		}
	}

	// resolve final values
	for i := range context.Args {
		option := context.Args[i]
		candidates := option.valueCandidates()
		if len(candidates) > 1 {
			return nil, &errorContext{
				internalErr: fmt.Errorf("multiple commands were provided for the same command. are you using different kinds of the same flag?"),
				errCode:     errParseOptions,
			}
		}
		if len(candidates) == 1 {
			option.Value = candidates[0]
		}
	}
	return context, nil
}

func (p *parser) performUsageAction(cmd *Command) (*Meta, *errorContext) {
	// at this point, if the command is interactive, execute its action and return the result
	if cmd.Interactive == true {
		args := p.getInteractiveArgs(cmd)
		return p.parseOptionsAndPerformAction(cmd, args)
	}
	// else there's nothing to do, so just display help
	return nil, p.output.printUsage(cmd)
}

func (p *parser) performAction(cmd *Command, ctx *Context) (*Meta, *errorContext) {
	if cmd.Action == nil {
		//panic("missing .Action for: " + cmd.Name)
		return nil, p.output.printUsage(cmd)
	}
	if p.isActionsEnabled {
		meta, err := cmd.Action(ctx)
		if err != nil {
			return meta, &errorContext{
				internalErr: err,
				errCode:     errCommandAction,
			}
		}
		return meta, nil
	}
	return nil, nil
}

func (p *parser) performErrorAction(cmd *Command, err *errorContext) *errorContext {
	if err.internalErr == flag.ErrHelp {
		return p.output.printUsage(cmd)
	}
	if err.errCode == errBashCompletion {
		return p.output.printCompletionForCommand(cmd)
	}
	return err
}

func (p *parser) commandForArg(commands []*Command, arg string) *Command {
	for _, cmd := range commands {
		if cmd.Name == arg {
			return cmd
		}
	}
	return nil
}

func (p *parser) getInteractiveArgs(cmd *Command) []string {
	args := []string{}
	for _, option := range cmd.Options {
		p.output.printOption(option)
		if option.IsSecure {
			args = append(args, "-"+option.Flag, ReadSecure())
		} else {
			args = append(args, "-"+option.Flag, Read())
		}
	}
	output.Println("") // force some space between interactive input and forthcoming output
	return args
}

func (p *parser) numOfMissingOptions(opts []*Option) int {
	missing := 0
	for _, o := range opts {
		if o.DefaultValue == nil {
			missing++
		}
	}
	return missing
}
