//go:build darwin

package execute

import (
	"fmt"
	"os/exec"
	"os/user"
	"strconv"
	"syscall"
)

// configureUser sets up the command to run as a different user on macOS.
func (e *BaseExecutor) configureUser(cmd *exec.Cmd) error {
	u, err := user.Lookup(e.user)
	if err != nil {
		return fmt.Errorf("failed to lookup user %q: %w", e.user, err)
	}

	uid, err := strconv.ParseUint(u.Uid, 10, 32)
	if err != nil {
		return fmt.Errorf("failed to parse uid: %w", err)
	}

	gid, err := strconv.ParseUint(u.Gid, 10, 32)
	if err != nil {
		return fmt.Errorf("failed to parse gid: %w", err)
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Credential: &syscall.Credential{
			Uid: uint32(uid),
			Gid: uint32(gid),
		},
	}

	logger.Debug("configured user", "user", e.user, "uid", uid, "gid", gid)
	return nil
}
