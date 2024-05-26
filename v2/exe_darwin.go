package execute

import (
	"os"
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
	return &DarwinExecutor{
		Environment: env,
		User:        user,
	}
}

// DarwinExecutor is an Executor implementation for Darwin systems.
type DarwinExecutor struct {
	BaseExecutor
	Environment []string
	User        string
}

// configureUser sets the user and group for the command to be executed.
func (e DarwinExecutor) configureUser(ctx context.Context, cancel context.CancelFunc, exe *exec.Cmd) error {
	u, err := user.Lookup(e.User)
	if err != nil {
		return nil, ctx, cancel, err
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return nil, ctx, cancel, err
	}

	gid, err := strconv.Atoi(u.Gid)
	if err != nil {
		return nil, ctx, cancel, err
	}

	exe.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)},
	}
}
