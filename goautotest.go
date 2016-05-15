package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/howeyc/fsnotify"
)

func startGoTest(doneChan chan bool) {
	fmt.Println("Running tests...")

	args := append([]string{"test"}, os.Args[1:]...)
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println()
	doneChan <- true
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = watcher.Watch(wd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
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
				go startGoTest(doneChan)
			}

		case err := <-watcher.Error:
			fmt.Println(err)

		case <-doneChan:
			running = false
		}
	}

}
