# Barbarian

Barbarian is a Go library which provides a convenient interface to run your program concurrently,
such as weak password blasting, port scanning, etc.

Inspired by [brutemachine](https://github.com/evilsocket/brutemachine), but I used channel to synchronzing goroutines.

## Installation

`go get github.com/AmyangXYZ/barbarian`

## Example

```go
package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/AmyangXYZ/barbarian"
)

func main() {
	bb := barbarian.New(DoRequest, OnResult, ReadFile("dict.txt"), 18)
	go func() {
		time.Sleep(3 * time.Second)
		bb.Stop()
	}()
	bb.Run()

	fmt.Printf("%+v\n", bb.Stats)
}

// DoRequest implements RunHandler.
func DoRequest(path string) interface{} {
	resp, err := http.Head("http://target.com" + path)
	// Only pass valid result to the handler.
	if err == nil && resp.StatusCode == 200 {
		return path
	}

	return nil
}

// OnResult implements ResHandler.
func OnResult(res interface{}) {
	fmt.Printf("@ Found '%s'\n", res)
}

// ReadFile reads file in lines.
func ReadFile(f string) (data []string) {
	b, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	defer b.Close()
	scanner := bufio.NewScanner(b)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}
	return
}
```