# go-execute

A Go library for safe and ergonomic execution of system commands.

## Design Principles

### Safety First

go-execute is designed to make it safer to execute system commands from Go applications, particularly in contexts where security matters—such as web applications, APIs, and tools that process user input.

**Key safety features:**

- **Direct execution by default**: Commands are executed directly without invoking a shell, eliminating an entire class of shell injection vulnerabilities. The executable is resolved via `exec.LookPath` and called directly with arguments passed as a list, not a string.

- **Explicit shell opt-in**: Shell execution is only used when explicitly requested via `WithShell()`. This makes dangerous patterns visible in code review and ensures developers consciously choose when shell features (pipes, redirects, globbing) are needed.

- **Structured command building**: Safe command wrappers provide a builder pattern API that constructs commands programmatically, avoiding string concatenation of user input.

### Usability

go-execute aims to be the most ergonomic way to execute commands in Go:

- **Simple things are simple**: Basic command execution requires minimal code
- **Complex things are possible**: Full control over environment, working directory, user context, timeouts, and async execution
- **Builder pattern for common commands**: Type-safe wrappers for grep, ls, cat, and find with structured results
- **Cross-platform**: Abstracts OS differences so the same code works on Linux, macOS, and Windows

## Installation

```bash
go get github.com/bgrewell/go-execute
```

## Usage

### Basic Execution (Safe by Default)

```go
// Direct execution - no shell involved
result, err := execute.Run("ls", "-la", "/tmp")
if err != nil {
    log.Fatal(err)
}
fmt.Println(result.Stdout)
```

### When You Need a Shell

```go
// Explicit opt-in for shell features (pipes, redirects, etc.)
exec := execute.NewExecutor(execute.WithShell("/bin/bash"))
result, err := exec.Execute("cat /var/log/*.log | grep error | head -10")
```

### Safe Command Wrappers

```go
// Type-safe builder pattern - immune to injection
result, err := execute.Grep(userPattern).
    InPath("/var/log").
    Recursive().
    IgnoreCase().
    Run()

for _, match := range result.Matches {
    fmt.Printf("%s:%d: %s\n", match.File, match.Line, match.Content)
}
```

### Async Execution

```go
result, err := execute.RunAsync("long-running-command", "arg1", "arg2")
if err != nil {
    log.Fatal(err)
}

// Stream output as it arrives
go io.Copy(os.Stdout, result.Stdout)
go io.Copy(os.Stderr, result.Stderr)

// Wait for completion
finalResult, err := result.Wait()
```

## Examples

See the [examples](./examples) directory for more usage examples.

## License

MIT
