
# Enhanced Go Error Group with Middleware Support

## Overview

`errgroup` is a lightweight Go library that extends [`golang.org/x/sync/errgroup`](https://pkg.go.dev/golang.org/x/sync/errgroup) by adding support for **middlewares**. These middlewares are executed **each time a function runs**, enabling additional functionality such as logging, metrics, or error handling.

## Installation

```sh
go get github.com/sysradium/errgroup
```

## Usage

### 1. Basic Usage (Same as errgroup.Group)

You can use Group just like errgroup.Group, but with middleware support.

```go
package main

import (
 "context"
 "fmt"
 "log"
 "log/slog"

 "github.com/sysradium/errgroup"
)

func main() {
 g := errgroup.New(errgroup.SlogLogger(slog.Default()))

 // Launch concurrent tasks
 g.Go(func() error {
  fmt.Println("Task 1")
  return nil
 })

 g.Go(func() error {
  fmt.Println("Task 2")
  return nil
 })

 // Wait for all tasks to finish
 if err := g.Wait(); err != nil {
  log.Fatalf("Error: %v", err)
 }
}
```

### 2. Using Middlewares

Middlewares are functions that wrap around task execution. They can be used for **logging, panic recovery, metrics, or any cross-cutting concern**.

#### Example: Logging Middleware with slog

The following middleware logs function calls before and after execution.

```go
package main

import (
 "context"
 "fmt"
 "log/slog"
 "os"

 "github.com/sysradium/errgroup"
)

func SlogLogger(l *slog.Logger) errgroup.Middleware {
 return func(next errgroup.GFn) errgroup.GFn {
  return func() error {
   err := next()
   if err != nil {
    l.Error("function call failed", "error", err)
   } else {
    l.Debug("function call succeeded")
   }
   return err
  }
 }
}

func main() {
 ctx := context.Background()
 logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

 // Create an errgroup with logging middleware
 g := errgroup.New(SlogLogger(logger))

 // Add tasks
 g.Go(func() error {
  return nil
 })

 g.Go(func() error {
  return fmt.Errorf("an error occurred")
 })

 // Wait for tasks to finish
 if err := g.Wait(); err != nil {
  fmt.Printf("error: %v\n", err)
 }
}
```

## API Reference

### New(middlewares ...Middleware) Group

Creates a new errgroup-middleware.Group instance with optional middlewares.

```go
g := errgroup.New(LoggingMiddleware1(), LoggingMiddleware2())
```

### WithContext(ctx context.Context) (Group, context.Context)

Returns an errgroup-middleware.Group tied to a context, just like errgroup.WithContext.

```go
g, ctx := errgroup.WithContext(context.Background())
```

### Go(fn GFn)

Runs a function asynchronously, applying middlewares.

```go
g.Go(func() error {
 fmt.Println("Running task")
 return nil
})
```

## License

This project is licensed under the MIT License.

## Contributing

Feel free to submit issues and PRs to improve the project!
