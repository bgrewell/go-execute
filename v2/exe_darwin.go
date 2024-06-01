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
	e := &DarwinExecutor{}
	for _, option := range options {
		option(e)
	}
	return e
}

// DarwinExecutor is an Executor implementation for Darwin systems.
type DarwinExecutor struct {
	BaseExecutor
}

// configureUser sets the user and group for the command to be executed.
func (e DarwinExecutor) configureUser(ctx context.Context, cancel context.CancelFunc, exe *exec.Cmd) error {
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
