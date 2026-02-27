# loginter

A Go linter compatible with [golangci-lint](https://golangci-lint.run/) that analyzes log messages and checks them against established rules.

## Rules

| # | Rule | Auto-fix | Description |
|---|------|----------|-------------|
| 1 | **Lowercase first letter** | Yes | Log messages must start with a lowercase letter |
| 2 | **English only** | No | Log messages must contain only English (ASCII) characters |
| 3 | **No special characters** | Yes | Log messages must not contain special characters or emoji |
| 4 | **No sensitive data** | No | Log messages must not contain potentially sensitive data |

### Supported loggers

- `log/slog` (standard library) - package-level functions and `*slog.Logger` methods
- `go.uber.org/zap` - `*zap.Logger` and `*zap.SugaredLogger` methods

### Examples

```go
// Rule 1: Lowercase first letter
slog.Info("Starting server")  // BAD
slog.Info("starting server")  // GOOD

// Rule 2: English only
slog.Info("запуск сервера")    // BAD
slog.Info("starting server")   // GOOD

// Rule 3: No special characters or emoji
slog.Info("server started! 🚀") // BAD
slog.Info("server started")     // GOOD

// Rule 4: No sensitive data
slog.Info("user password: " + password)  // BAD
slog.Info("user authenticated")          // GOOD
```

## Installation

### As a golangci-lint plugin (recommended)

1. Create `.custom-gcl.yml` in your project root:

```yaml
version: v2.2.0
plugins:
  - module: 'github.com/anisimov-anthony/loginter'
    import: 'github.com/anisimov-anthony/loginter/plugin'
    version: v1.0.0
```

2. Add to your `.golangci.yml`:

```yaml
version: "2"
linters:
  enable:
    - loginter
  settings:
    custom:
      loginter:
        type: "module"
        description: "Checks log messages for common issues"
        settings:
          check_lowercase: true
          check_english: true
          check_special: true
          check_sensitive: true
          sensitive_patterns:
            - "ssn"
            - "credit_card"
```

3. Build and run:

```bash
golangci-lint custom
./custom-gcl run ./...
```

4. Auto-fix (where supported):

```bash
./custom-gcl run --fix ./...
```

### As a standalone tool

```bash
go install github.com/anisimov-anthony/loginter/cmd/loginter@latest
loginter ./...
```

## Configuration

All checks are enabled by default. You can disable individual checks via the golangci-lint settings:

| Setting | Type | Default | Description |
|---------|------|---------|-------------|
| `check_lowercase` | `bool` | `true` | Check that log messages start with a lowercase letter |
| `check_english` | `bool` | `true` | Check that log messages are in English only |
| `check_special` | `bool` | `true` | Check for special characters and emoji |
| `check_sensitive` | `bool` | `true` | Check for sensitive data patterns |
| `sensitive_patterns` | `[]string` | `[]` | Additional sensitive data patterns (added to defaults) |

### Default sensitive data patterns

The following patterns are checked by default: `password`, `passwd`, `secret`, `token`, `api_key`, `apikey`, `api_secret`, `access_key`, `private_key`, `credential`, `auth`.

Pattern matching respects word boundaries, so `"user authenticated"` will **not** trigger a match for `auth`, but `"auth header set"` will.

You can add your own patterns via the `sensitive_patterns` setting.

## Building from source

```bash
git clone https://github.com/anisimov-anthony/loginter.git
cd loginter
go build ./...
```

## Running tests

### Unit tests

```bash
go test ./... -v -cover
```

### Integration (e2e) tests

Integration tests build the standalone `loginter` binary and run it against real Go source files, verifying exit codes and diagnostic output.

```bash
go test -tags e2e ./e2e/ -v
```

To run everything at once:

```bash
go test ./... && go test -tags e2e ./e2e/ -v
```

The e2e tests are kept behind the `e2e` build tag so they don't slow down the normal `go test ./...` cycle. They are run automatically in CI as a separate step.


## License

`MIT`
