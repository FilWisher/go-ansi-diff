# go-ansi-diff
a literal translation of https://github.com/mafintosh/ansi-diff-stream to go

## install

```go get github.com/filwisher/go-ansi-diff```

## example

```go
package main

import (
  diff "github.com/filwisher/go-ansi-diff"
  "time"
  "fmt"
)

func main() {

	d := diff.Differ{}

	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Printf("%s", d.Diff([]byte(fmt.Sprintf("%s", time.Now().String()))))
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	<-quit
}

```
