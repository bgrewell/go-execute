package execute

import "runtime"

type Option func(executor Executor)

func WithEnvironment(env []string) Option {
	return func(e Executor) {
		e.SetEnvironment(env)
	}
}

func WithUser(user string) Option {
	return func(e Executor) {
		e.SetUser(user)
	}
}

func WithDefaultShell() Option {
	shell := "/bin/bash"
	if runtime.GOOS == "windows" {
		shell = "cmd.exe"
	}
	return func(e Executor) {
		e.SetShell(shell)
	}
}

func WithShell(shell string) Option {
	return func(e Executor) {
		e.SetShell(shell)
	}
}

func WithWorkingDir(dir string) Option {
	return func(e Executor) {
		e.SetWorkingDir(dir)
	}
}

func WithSudoCredentials(password string) Option {
	return func(e Executor) {
		e.SetSudoCredentials(password)
	}
}
