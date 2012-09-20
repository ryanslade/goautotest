package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.Watch(wd)
	if err != nil {
		log.Fatal(err)
	}

	defer watcher.Close()

	for {
		select {
		case ev := <-watcher.Event:
			if strings.HasSuffix(ev.Name, ".go") {
				cmd := exec.Command("go", "test")
				stdout, err := cmd.StdoutPipe()
				if err != nil {
					log.Println(err)
					break
				}

				err = cmd.Start()
				if err != nil {
					log.Println(err)
					break
				}

				go io.Copy(os.Stdout, stdout)
				err = cmd.Wait()
				if err != nil {
					log.Println(err)
				}

				fmt.Println()
			}

		case err := <-watcher.Error:
			log.Println(err)
		}
	}

}
