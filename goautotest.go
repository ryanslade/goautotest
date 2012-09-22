package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

func startGoTest() {
	fmt.Println("Running tests...")

	args := append([]string{"test"}, os.Args[1:]...)
	cmd := exec.Command("go", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go io.Copy(os.Stdout, stdout)
	err = cmd.Wait()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println()
}

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = watcher.Watch(wd)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer watcher.Close()

	burstGuard := make(<-chan time.Time)
	for {
		select {
		case ev := <-watcher.Event:
			if strings.HasSuffix(ev.Name, ".go") {
				burstGuard = time.After(500 * time.Millisecond)
			}

		case err := <-watcher.Error:
			fmt.Println(err)

		case <-burstGuard:
			startGoTest()
		}
	}

}
