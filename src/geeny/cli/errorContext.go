package cli

type errorContext struct {
	internalErr error
	errCode     errorType
}

type errorType int

const (
	errParseOptions errorType = iota
	errCommandNotSupported
	errCommandAction
	errUsage
	errBashCompletion
)
