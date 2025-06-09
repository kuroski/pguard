# Graceful Runner

A small Go utility that runs a command with optional timeout and graceful shutdown handling.

## Usage

```bash
go run main.go [flags] <command> [args...]
```

### Flags

| Flag        | Type     | Description                                                                  |
| ----------- | -------- | ---------------------------------------------------------------------------- |
| `-timeout`  | duration | (Optional) Maximum allowed time for the command to run. e.g. `30s`, `2m`.    |
| `-graceful` | duration | (Optional) Time to wait after `SIGINT` before killing the process. e.g. `5s` |

_Note_: Use negative values to disable timeout or graceful shutdown respectively.

### Example

```bash
pguard -timeout=10s -graceful=3s sleep 100
```

This will:

- Run `sleep 100`
- Timeout after 10 seconds
- If interrupted (`Ctrl+C`), give it 3 seconds to exit before forcefully killing it
