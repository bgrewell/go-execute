# go-execute

Safe, ergonomic command execution for Go.

```go
result, err := execute.Run("git", "status", "--porcelain")
```

## Why go-execute?

Go's `os/exec` is powerful but makes it easy to write code vulnerable to shell injection. go-execute provides a safer default while remaining simple to use.

| Feature | go-execute | os/exec |
|---------|-----------|---------|
| Shell injection protection | Default | Manual |
| Structured output parsing | Built-in | DIY |
| Async with proper pipe handling | Built-in | Tricky |
| Cross-platform command wrappers | Yes | No |

## Installation

```bash
go get github.com/bgrewell/go-execute
```

## Quick Start

```go
import "github.com/bgrewell/go-execute"

// Simple command
result, err := execute.Run("ls", "-la")
fmt.Println(result.Stdout)

// With options
exec := execute.NewExecutor(
    execute.WithWorkingDir("/tmp"),
    execute.WithEnvironment([]string{"DEBUG=1"}),
)
result, err := exec.Execute("make", "build")

// Async execution
async, _ := execute.RunAsync("long-task")
go io.Copy(os.Stdout, async.Stdout)
result, err := async.Wait()
```

## Safe by Default

Commands execute **directly without a shell**, preventing injection attacks:

```go
// Safe: argument treated as literal, not interpreted
execute.Run("echo", userInput)  // userInput = "hello; rm -rf /"
// Runs: echo "hello; rm -rf /"

// Shell required? Opt-in explicitly:
exec := execute.NewExecutor(execute.WithShell("/bin/bash"))
exec.Execute("ls *.go | head -5")
```

## Command Wrappers

Type-safe builders for common commands with structured results:

```go
// Grep with parsed results
result, _ := execute.Grep("TODO").
    InPath("./src").
    Recursive().
    IgnoreCase().
    Run()

for _, m := range result.Matches {
    fmt.Printf("%s:%d: %s\n", m.File, m.Line, m.Content)
}

// List files
files, _ := execute.Ls("/var/log").Long().Run()
for _, f := range files.Entries {
    fmt.Printf("%s %d bytes\n", f.Name, f.Size)
}
```

Available wrappers: `Grep`, `Ls`, `Cat`, `Find`

## Configuration

```go
exec := execute.NewExecutor(
    execute.WithWorkingDir("/path/to/dir"),
    execute.WithEnvironment([]string{"KEY=value"}),
    execute.WithShell("/bin/bash"),      // Enable shell mode
    execute.WithUser("www-data"),        // Run as user (Unix)
    execute.WithSudoPassword("..."),     // For sudo commands
)
```

## Context Support

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := execute.RunContext(ctx, "slow-command")
if errors.Is(err, context.DeadlineExceeded) {
    // Handle timeout
}
```

## Examples

See the [examples](./examples) directory for complete working examples.

## License

MIT
