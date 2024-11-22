package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/howeyc/fsnotify"
)

var args = append([]string{"test"}, os.Args[1:]...)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fail(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		fail(err)
	}

	if err := watcher.Watch(wd); err != nil {
		fail(err)
	}
	defer watcher.Close()

	running := false
	doneChan := make(chan bool)

	for {
		select {
		case ev := <-watcher.Event:
			if running {
				continue
			}
			if strings.HasSuffix(ev.Name, ".go") {
				running = true
				go test(doneChan)
			}

		case err := <-watcher.Error:
			fmt.Println(err)

		case <-doneChan:
			running = false
		}
	}
}

func test(doneChan chan bool) {
	fmt.Println("Running tests...")

	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fail(err)
	}

	fmt.Println()
	doneChan <- true
}

func fail(err error) {
	fmt.Println(err)
	os.Exit(1)
}
