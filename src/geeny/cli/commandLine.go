package cli

// FlagGenerateBashCompletion can be sent into the cli to generate bash completion subcommands / options for a given command
const FlagGenerateBashCompletion = "generate-bash-completion"

// CommandLine encapsulates command line functionality
type CommandLine struct {
	parser *parser
}

// NewCommandLine create a new command line instance
func NewCommandLine() *CommandLine {
	return &CommandLine{
		parser: newParser(true),
	}
}

// Run starts the cli
func (cl *CommandLine) Run(command *Command, args []string) (*Meta, error) {
	meta, err := cl.parser.parseCommand(command, args, 0)
	if err != nil {
		return meta, err.internalErr
	}
	return meta, nil
}

// Search searches an array of arguments for a given flag
func (cl *CommandLine) Search(args []string, flag string) bool {
	for _, arg := range args {
		if arg == "-"+flag || arg == "--"+flag {
			return true
		}
	}
	return false
}

// SetActionsEnabled disables actions
func (cl *CommandLine) SetActionsEnabled(enabled bool) {
	cl.parser.isActionsEnabled = enabled
}

// GetActionsEnabled returns true if actions are enabled
func (cl *CommandLine) GetActionsEnabled() bool {
	return cl.parser.isActionsEnabled
}
