package execute

import (
	"context"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

// NewExecutor creates a new Executor.
func NewExecutor(options ...Option) Executor {
	e := &LinuxExecutor{}
	for _, option := range options {
		option(e)
	}
	return e
}

// LinuxExecutor is an Executor implementation for Linux systems.
type LinuxExecutor struct {
	BaseExecutor
}

// configureUser sets the user and group for the command to be executed.
func (e LinuxExecutor) configureUser(ctx context.Context, cancel context.CancelFunc, exe *exec.Cmd) error {
	u, err := user.Lookup(e.user)
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
