Package vlog implements a simple logging package on top of existing Go log package.
vlog package adds logging levels support, otherwise vlog package behaves exactly as the Go log package.
vlog does not create any hirarchy of the loggers. Root logger is always created with name "".

# Quick start

```golang
logger := vlog.GetLogger("example", nil)

logger.Tracef("%s, a trace log", "Hello!")
logger.Debug("Hello!, a debug log")

// Output:
// trace Hello!, a trace log
// debug Hello!, a debug log
```

# Configuration

Set logging level of the `Logger` using `SetLevel` method.

```golang
logger := vlog.GetLogger("example", nil)
logger.SetLevel(vlog.INFO)
```

Set flags of the `Logger` using `SetFlags` method.

```golang
logger := vlog.GetLogger("example", nil)
logger.SetFlags(log.Ldate | log.Ltime)
```
