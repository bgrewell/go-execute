//go:build windows

package execute

import (
	"fmt"
	"os/exec"
)

// configureUser sets up the command to run as a different user on Windows.
// Note: Running as a different user on Windows requires different approaches
// such as runas or token manipulation, which is not yet implemented.
func (e *BaseExecutor) configureUser(cmd *exec.Cmd) error {
	// TODO: Implement Windows user impersonation
	// Options include:
	// - Using runas command wrapper
	// - Token manipulation via Windows API
	// - CreateProcessWithLogonW
	return fmt.Errorf("user impersonation not yet implemented on Windows")
}
