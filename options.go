package execute

import "runtime"

// Option is a functional option for configuring a BaseExecutor
type Option func(*BaseExecutor)

// WithEnvironment sets the environment variables for command execution.
// The format is the same as os.Environ(): "KEY=VALUE" strings.
func WithEnvironment(env []string) Option {
	return func(e *BaseExecutor) {
		e.env = env
	}
}

// WithWorkingDir sets the working directory for command execution
func WithWorkingDir(dir string) Option {
	return func(e *BaseExecutor) {
		e.workingDir = dir
	}
}

// WithShell sets a custom shell for command execution
func WithShell(shell string) Option {
	return func(e *BaseExecutor) {
		e.shell = shell
	}
}

// WithDefaultShell sets the OS-appropriate default shell.
// On Unix systems this is "/bin/bash", on Windows it is "cmd.exe".
func WithDefaultShell() Option {
	return func(e *BaseExecutor) {
		if runtime.GOOS == "windows" {
			e.shell = "cmd.exe"
		} else {
			e.shell = "/bin/bash"
		}
	}
}

// WithUser sets the user to run the command as
func WithUser(user string) Option {
	return func(e *BaseExecutor) {
		e.user = user
	}
}

// WithSudoPassword sets the password for sudo operations.
// The password is stored securely in memory.
func WithSudoPassword(password string) Option {
	return func(e *BaseExecutor) {
		e.sudoPassword = password
	}
}
