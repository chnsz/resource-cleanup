package helper

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func Run(cmd *exec.Cmd, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Error creating stderr pipe:", err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading from stdout:", err)
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading from stderr:", err)
		}
	}()

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Command execution timed out")
		_ = cmd.Process.Kill()
	case err := <-done:
		if err != nil {
			fmt.Println("Error waiting for command:", err)
		}
	}
}

func WriteToFile(path, content string) error {
	path = strings.ReplaceAll(path, "\\", "/")

	pos := strings.LastIndex(path, "/")
	directory := path[:pos]

	if err := os.MkdirAll(directory, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %s, error: %s", path, err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %s, error: %s", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("failed to close file: %s", err)
		}
	}()

	if _, err = f.WriteString(content); err != nil {
		return fmt.Errorf("failed to write data to %s, error: %s", path, err)
	}
	return nil
}

func GetTmpDir() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 32)
	for i := range b {
		b[i] = charset[rnd.Intn(len(charset))]
	}
	return fmt.Sprintf("%s%s%s", os.TempDir(), string(os.PathSeparator), string(b))
}
