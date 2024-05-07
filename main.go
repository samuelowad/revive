package main

import (
	"github.com/samuelowad/revive/config"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
)

var (
	appCmd *exec.Cmd
	mutex  sync.Mutex
)

func main() {
	config.ReadConfig()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	if err := startCommand(); err != nil {
		log.Fatalf("Failed to start the command: %v", err)
	}

	if err := watchFiles(watcher); err != nil {
		log.Fatalf("Error watching files: %v", err)
	}
}

func startCommand() error {
	log.Println("Starting the command...")
	command := strings.Fields(config.Config.Command)
	appCmd = exec.Command(command[0], command[1:]...)
	appCmd.Stdout = os.Stdout
	appCmd.Stderr = os.Stderr
	return appCmd.Start()
}

func restartCommand() error {
	log.Println("Restarting the command...")
	mutex.Lock()
	defer mutex.Unlock()

	if appCmd != nil && appCmd.Process != nil {
		if err := appCmd.Process.Signal(syscall.SIGTERM); err != nil {
			return err
		}
		if err := appCmd.Wait(); err != nil && !strings.Contains(err.Error(), "signal: terminated") {
			return err
		}
	}

	if config.Config.RestartDelaySeconds > 0 {
		time.Sleep(time.Duration(config.Config.RestartDelaySeconds) * time.Second)
	}

	return startCommand()
}

func watchFiles(watcher *fsnotify.Watcher) error {
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
		return err
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Printf("File %s: %s\n", actionName(event.Op), event.Name)
				if shouldMonitorFile(event.Name) {
					restartCommand()
				}
			}
		case err := <-watcher.Errors:
			log.Println("Error:", err)
		}
	}
}

func shouldMonitorFile(filename string) bool {
	dir := filepath.Dir(filename)
	for _, ignoredDir := range config.Config.IgnoreDirectories {
		if dir == ignoredDir {
			return false
		}
	}

	ext := filepath.Ext(filename)
	for _, ignoredEnding := range config.Config.IgnoreFileNameEndsWith {
		if strings.HasSuffix(filename, ignoredEnding) {
			return false
		}
	}

	for _, monitoredExt := range config.Config.MonitorFileExt {
		if ext == monitoredExt {
			return true
		}
	}

	return false
}

func actionName(op fsnotify.Op) string {
	switch {
	case op&fsnotify.Write == fsnotify.Write:
		return "Modified"
	case op&fsnotify.Create == fsnotify.Create:
		return "Created"
	case op&fsnotify.Remove == fsnotify.Remove:
		return "Deleted"
	default:
		return "Unknown"
	}
}
