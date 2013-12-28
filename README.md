# Daemon

A daemon library for Go.

## Install

```
$ go get "github.com/yoppi/daemon"
```

## How to use

Call `daemon#Daemonize()`, then starting on background.

```go
import (
  "github.com/yoppi/daemon"
)

func main() {
  daemon.Daemonize("app.log")
  // start background process
  for {
    fmt.Println("Background ....")
  }
}
```

