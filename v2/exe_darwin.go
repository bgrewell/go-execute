package v2

func NewExecutor(env []string) Executor {
	return NewExecutorAsUser("", env)
}

func NewExecutorAsUser(user string, env []string) Executor {
	return &LinuxExecutor{
		Environment: env,
		User:        user,
	}
}

type DarwinExecutor struct {
	Environment []string
	User        string
}

func (e DarwinExecutor) Execute(command string) (combined string, err error) {
	//TODO implement me
	panic("implement me")
}

func (e DarwinExecutor) ExecuteSeparate(command string) (stdout string, stderr string, err error) {
	//TODO implement me
	panic("implement me")
}

func (e DarwinExecutor) ExecuteStream(command string) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	//TODO implement me
	panic("implement me")
}

func (e DarwinExecutor) ExecuteStreamWithInput(command string, stdin io.WriteCloser) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	//TODO implement me
	panic("implement me")
}

func (e DarwinExecutor) ExecuteWithTimeout(command string, timeout time.Duration) (combined string, err error) {
	//TODO implement me
	panic("implement me")
}

func (e DarwinExecutor) ExecuteSeparateWithTimeout(command string, timeout time.Duration) (stdout string, stderr string, err error) {
	//TODO implement me
	panic("implement me")
}

func (e DarwinExecutor) ExecuteStreamWithTimeout(command string, timeout time.Duration) (stdout io.ReadCloser, stderr io.ReadCloser, err error) {
	//TODO implement me
	panic("implement me")
}
