package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/fsnotify/fsnotify"
)

var (
	commandName string
	commandArgs []string
	appCmd      *exec.Cmd
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: GoRevive <command> [arguments...]")
	}
	commandName = os.Args[1]
	commandArgs = os.Args[2:]

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if err := startCommand(); err != nil {
		log.Fatalf("Failed to start the command: %v", err)
	}

	go watchFiles(watcher)

	done := make(chan bool)
	<-done
}

func startCommand() error {
	log.Println("Starting the command...")
	appCmd = exec.Command(commandName, commandArgs...)
	appCmd.Stdout = os.Stdout
	appCmd.Stderr = os.Stderr
	return appCmd.Start()
}

func restartCommand() error {
	log.Println("Restarting the command...")
	err := appCmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}
	err = appCmd.Wait()
	if err != nil && !strings.Contains(err.Error(), "signal: terminated") {
		return err
	}
	return startCommand()
}

func watchFiles(watcher *fsnotify.Watcher) {
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event := <-watcher.Events:
			switch {
			case event.Op&fsnotify.Write == fsnotify.Write:
				handleEvent(event, "Modified")
			case event.Op&fsnotify.Create == fsnotify.Create:
				handleEvent(event, "Created")
			case event.Op&fsnotify.Remove == fsnotify.Remove:
				handleEvent(event, "Deleted")
			}
		case err := <-watcher.Errors:
			log.Println("Error:", err)
		}
	}
}

func handleEvent(event fsnotify.Event, action string) {
	if strings.HasSuffix(event.Name, ".go") {
		log.Printf("File %s: %s\n", action, event.Name)
		restartCommand()
	}
}
