package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"io"
	"os"
	"os/exec"
	"strings"
)

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

	for {
		select {
		case ev := <-watcher.Event:
			if strings.HasSuffix(ev.Name, ".go") {
				fmt.Println("Running tests...")

				args := append([]string{"test"}, os.Args[1:]...)
				cmd := exec.Command("go", args...)

				stdout, err := cmd.StdoutPipe()
				if err != nil {
					fmt.Println(err)
					break
				}

				err = cmd.Start()
				if err != nil {
					fmt.Println(err)
					break
				}

				go io.Copy(os.Stdout, stdout)
				err = cmd.Wait()
				if err != nil {
					fmt.Println(err)
				}

				fmt.Println()
			}

		case err := <-watcher.Error:
			fmt.Println(err)
		}
	}

}
