/*
Package logging contains logger adapters that make it possible for shogoa to log messages to various
logger backends. Each adapter exists in its own sub-package named after the corresponding logger
package.

Once instantiated adapters can be used by setting the shogoa service logger with WithLogger:

```go

	  func main() {
	    // ...

	    // Setup logger adapter
	    logger := log15.New()

	    // Create service
	    service := shogoa.New("my service")
	    service.WithLogger(goalog15.New(logger))

	    // ...
	}

```

See http://shogoa.design/implement/logging/ for details.
*/
package logging