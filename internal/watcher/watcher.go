package watcher

import (
	"bufio"
	"os"
	"time"
)

// Watch watches a file and sends new lines on the channel
func Watch(path string, lines chan<- string, stop <-chan struct{}) error {
	// Wait for file to exist
	for {
		if _, err := os.Stat(path); err == nil {
			break
		}

		select {
		case <-stop:
			return nil
		case <-time.After(1 * time.Second):
			// Retry
		}
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = file.Close()
	}()

	// Seek to end
	_, _ = file.Seek(0, 2)

	scanner := bufio.NewScanner(file)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			return nil
		case <-ticker.C:
			for scanner.Scan() {
				lines <- scanner.Text()
			}
		}
	}
}
