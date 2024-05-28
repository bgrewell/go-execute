package execute

import (
	"reflect"
	"testing"
)

func TestSetEnvironment(t *testing.T) {
	executor := &BaseExecutor{}

	t.Run("SetEnvironment_WithValidEnvironment", func(t *testing.T) {
		env := []string{"VAR1=value1", "VAR2=value2"}
		executor.SetEnvironment(env)

		if len(executor.Environment()) != len(env) {
			t.Errorf("Expected environment length %d, but got %d", len(env), len(executor.Environment()))
		}
	})

	t.Run("SetEnvironment_WithEmptyEnvironment", func(t *testing.T) {
		env := []string{}
		executor.SetEnvironment(env)

		if len(executor.Environment()) != len(env) {
			t.Errorf("Expected environment length %d, but got %d", len(env), len(executor.Environment()))
		}
	})

	t.Run("SetEnvironment_WithNilEnvironment", func(t *testing.T) {
		executor.SetEnvironment(nil)

		if executor.Environment() != nil {
			t.Errorf("Expected environment to be nil, but got %v", executor.Environment())
		}
	})
}

func TestEnvironment(t *testing.T) {
	executor := &BaseExecutor{}

	t.Run("Environment_WhenEnvironmentIsSet", func(t *testing.T) {
		env := []string{"VAR1=value1", "VAR2=value2"}
		executor.SetEnvironment(env)

		if !reflect.DeepEqual(executor.Environment(), env) {
			t.Errorf("Expected environment %v, but got %v", env, executor.Environment())
		}
	})

	t.Run("Environment_WhenEnvironmentIsNotSet", func(t *testing.T) {
		executor := &BaseExecutor{}

		if len(executor.Environment()) != 0 {
			t.Errorf("Expected environment to be empty, but got %v", executor.Environment())
		}
	})
}

func TestSetUser(t *testing.T) {
	executor := &BaseExecutor{}

	t.Run("SetUser_WithValidUser", func(t *testing.T) {
		user := "testUser"
		executor.SetUser(user)

		if executor.User() != user {
			t.Errorf("Expected user %s, but got %s", user, executor.User())
		}
	})

	t.Run("SetUser_WithEmptyUser", func(t *testing.T) {
		user := ""
		executor.SetUser(user)

		if executor.User() != user {
			t.Errorf("Expected user to be empty, but got %s", executor.User())
		}
	})
}

func TestUser(t *testing.T) {
	executor := &BaseExecutor{}

	t.Run("User_WhenUserIsSet", func(t *testing.T) {
		user := "testUser"
		executor.SetUser(user)

		if executor.User() != user {
			t.Errorf("Expected user %s, but got %s", user, executor.User())
		}
	})

	t.Run("User_WhenUserIsNotSet", func(t *testing.T) {
		executor := &BaseExecutor{}

		if executor.User() != "" {
			t.Errorf("Expected user to be empty, but got %s", executor.User())
		}
	})
}

func TestSetShell(t *testing.T) {
	executor := &BaseExecutor{}

	t.Run("SetShell_WithValidShell", func(t *testing.T) {
		shell := "/bin/bash"
		executor.SetShell(shell)

		if executor.Shell() != shell {
			t.Errorf("Expected shell %s, but got %s", shell, executor.Shell())
		}
	})

	t.Run("SetShell_WithEmptyShell", func(t *testing.T) {
		shell := ""
		executor.SetShell(shell)

		if executor.Shell() != shell {
			t.Errorf("Expected shell to be empty, but got %s", executor.Shell())
		}
	})
}

func TestClearShell(t *testing.T) {
	executor := &BaseExecutor{}

	t.Run("ClearShell_WhenShellIsSet", func(t *testing.T) {
		shell := "/bin/bash"
		executor.SetShell(shell)
		executor.ClearShell()

		if executor.Shell() != "" {
			t.Errorf("Expected shell to be empty, but got %s", executor.Shell())
		}
	})

	t.Run("ClearShell_WhenShellIsAlreadyEmpty", func(t *testing.T) {
		executor.ClearShell()

		if executor.Shell() != "" {
			t.Errorf("Expected shell to be empty, but got %s", executor.Shell())
		}
	})
}

func TestShell(t *testing.T) {
	executor := &BaseExecutor{}

	t.Run("Shell_WhenShellIsSet", func(t *testing.T) {
		shell := "/bin/bash"
		executor.SetShell(shell)

		if executor.Shell() != shell {
			t.Errorf("Expected shell %s, but got %s", shell, executor.Shell())
		}
	})

	t.Run("Shell_WhenShellIsNotSet", func(t *testing.T) {
		executor := &BaseExecutor{}

		if executor.Shell() != "" {
			t.Errorf("Expected shell to be empty, but got %s", executor.Shell())
		}
	})
}

func TestUsingShell(t *testing.T) {
	executor := &BaseExecutor{}

	t.Run("UsingShell_WhenShellIsSet", func(t *testing.T) {
		shell := "/bin/bash"
		executor.SetShell(shell)

		if !executor.UsingShell() {
			t.Errorf("Expected UsingShell to be true, but got false")
		}
	})

	t.Run("UsingShell_WhenShellIsNotSet", func(t *testing.T) {
		executor := &BaseExecutor{}

		if executor.UsingShell() {
			t.Errorf("Expected UsingShell to be false, but got true")
		}
	})
}
