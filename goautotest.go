package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

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

	err = filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			err = watcher.Watch(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer watcher.Close()

	ignore := false
	doneChan := make(chan bool)
	readyChan := make(chan bool)

	for {
		select {
		case ev := <-watcher.Event:
			if strings.HasSuffix(ev.Name, ".go") && !ignore {
				ignore = true
				go startGoTest(doneChan)
			}

		case err := <-watcher.Error:
			fmt.Println(err)

		case <-doneChan:
			time.AfterFunc(1500*time.Millisecond, func() {
				readyChan <- true
			})

		case <-readyChan:
			ignore = false
		}
	}

}
