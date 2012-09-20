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

	err = watcher.Watch(".")
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case ev := <-watcher.Event:
			if strings.HasSuffix(ev.Name, ".go") {
				cmd := exec.Command("go", "test")
				stdout, err := cmd.StdoutPipe()
				if err != nil {
					log.Println("Error: ", err)
					break
				}

				err = cmd.Start()
				if err != nil {
					log.Println("Error running test: ", err)
					break
				}

				go io.Copy(os.Stdout, stdout)
				err = cmd.Wait()
				if err != nil {
					log.Println("Error:", err)
				}

				fmt.Println()
			}

		case err := <-watcher.Error:
			log.Println("Error:", err)
		}
	}

	watcher.Close()
}
