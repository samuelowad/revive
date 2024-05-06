package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGoRevive(t *testing.T) {
	// Create a temporary directory
	tmpDir, err := ioutil.TempDir("", "test-data")
	if err != nil {
		t.Fatal("Failed to create temporary directory:", err)
	}
	defer os.RemoveAll(tmpDir)

	// Copy test.go to the temporary directory
	testGoSrc := "test-data/test.go"
	testGoDst := filepath.Join(tmpDir, "test.go")
	if err := copyFile(testGoDst, testGoSrc); err != nil {
		t.Fatal("Failed to copy test.go:", err)
	}

	// Build GoRevive
	if err := buildGoRevive(tmpDir); err != nil {
		t.Fatal("Failed to build GoRevive:", err)
	}

	// Start GoRevive in the temporary directory
	goRevivePath := filepath.Join(tmpDir, "GoRevive")
	cmd := exec.Command(goRevivePath, "go", "run", "test.go")
	cmd.Dir = tmpDir
	if err := cmd.Start(); err != nil {
		t.Fatal("Failed to start GoRevive:", err)
	}
	defer cmd.Process.Kill()

	// Wait for the GoRevive script to start
	time.Sleep(2 * time.Second)

	// Modify test.go
	if err := modifyFile(testGoDst); err != nil {
		t.Fatal("Failed to modify test.go:", err)
	}

	// Wait for the script to reload the application
	time.Sleep(2 * time.Second)

	// Restore test.go
	if err := restoreFile(testGoDst, testGoSrc); err != nil {
		t.Fatal("Failed to restore test.go:", err)
	}
}

func buildGoRevive(tmpDir string) error {
	cmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "GoRevive"), "main.go")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func copyFile(dst, src string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func modifyFile(filename string) error {
	// Read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// Modify the content
	modifiedData := []byte(strings.Replace(string(data), "fmt.Println(\"Hello, playground\")", "fmt.Println(\"Modified!\")", 1))
	return os.WriteFile(filename, modifiedData, 0644)
}

func restoreFile(dst, src string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}
