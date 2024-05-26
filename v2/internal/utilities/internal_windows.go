package utilities

import (
	"errors"
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"golang.org/x/sys/windows"
	"os"
	"syscall"
)

// GetTokenForUser gets the token for the specified user.
func GetTokenForUser(user string) (syscall.Token, error) {
	// Try to find a process running as the target user
	pid := int32(0)
	processes, _ := process.Processes()
	for _, process := range processes {
		username, _ := process.Username()
		if username == user {
			pid = process.Pid
			break
		}
	}
	if pid == 0 {
		return 0, errors.New("unable to find process running as target user")
	}

	token, err := GetTokenFromPid(pid)
	if err != nil {
		return 0, err
	}

	return token, nil
}

// GetTokenFromPid gets the token for the specified process ID.
func GetTokenFromPid(pid int32) (syscall.Token, error) {
	var err error
	var token syscall.Token

	handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		fmt.Println("Token Process", "err", err)
		return 0, err
	}
	defer syscall.CloseHandle(handle)

	// Find process token via win32
	err = syscall.OpenProcessToken(handle, syscall.TOKEN_ALL_ACCESS, &token)
	if err != nil {
		fmt.Println("Open Token Process", "err", err)
		return 0, err
	}

	return token, nil
}

// RunningAsAdmin checks to see if the current process is running as an administrator by attempting to open a handle
// to the physical drive which typically requires admin privileges.
func RunningAsAdmin() bool {
	// Check to see if we are running with the right permissions
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

// HasRequiredPrivileges checks to see if the current process has the required privileges.
func HasRequiredPrivileges() (admin bool, elevated bool, err error) {
	var sid *windows.SID

	err = windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		return false, false, err
	}
	defer windows.FreeSid(sid)

	// Get the token for the active thread
	token := windows.Token(0)

	isAdmin, err := token.IsMember(sid)
	if err != nil {
		return false, false, err
	}

	return isAdmin, token.IsElevated(), nil
}

// LookupAccount looks up the account for the specified username.
func LookupAccount(username string) (*syscall.SID, error) {
	// Use the LookupAccountName syscall to verify the user exists
	var sid *syscall.SID
	var domain uint16
	var size uint32
	var peUse uint32
	err := syscall.LookupAccountName(nil, syscall.StringToUTF16Ptr(username), sid, &size, &domain, &size, &peUse)
	return sid, err
}
