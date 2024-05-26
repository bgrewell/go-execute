package execute

import (
	"context"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

// NewExecutor creates a new Executor.
func NewExecutor() Executor {
	return NewExecutorAsUser("", os.Environ())
}

// NewExecutorWithEnv creates a new Executor with the specified environment.
func NewExecutorWithEnv(env []string) Executor {
	return NewExecutorAsUser("", env)
}

// NewExecutorAsUser creates a new Executor with the specified user and environment.
func NewExecutorAsUser(user string, env []string) Executor {
	return &LinuxExecutor{
		Environment: env,
		User:        user,
	}
}

// LinuxExecutor is an Executor implementation for Linux systems.
type LinuxExecutor struct {
	BaseExecutor
	Environment []string
	User        string
}

// configureUser sets the user and group for the command to be executed.
func (e LinuxExecutor) configureUser(ctx context.Context, cancel context.CancelFunc, exe *exec.Cmd) error {
	u, err := user.Lookup(e.User)
	if err != nil {
		return err
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return err
	}

	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return err
	}

	exe.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)},
	}

	return nil
}
