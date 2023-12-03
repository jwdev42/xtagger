package global

const (
	ExitSuccess ProgramExitCode = iota
	ExitHardError
	ExitSoftError
)

const BufSize = 1048576 //Default buffer size is 1 MiB

type ProgramExitCode int

var ExitCode ProgramExitCode = ExitSuccess
