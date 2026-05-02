package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var progressRegex = regexp.MustCompile(`\[download\]\s+([0-9.]+)%`)

func downloadAudio(url string, dest string, quiet bool) error {
	cmd := exec.Command(
		"yt-dlp",
		"--extract-audio",
		"--audio-format", "mp3",
		"--no-warnings",
		"--newline",
		url,
	)

	cmd.Dir = dest

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %w", err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()

		matches := progressRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			percent := matches[1]

			if !quiet {
				fmt.Printf("\r\033[KDownloading: %s%%", percent)
			}
		}
	}

	if !quiet {
		fmt.Println()
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading output: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}

func worker(id int, jobs <-chan string, dest string, wg *sync.WaitGroup) {
	defer wg.Done()

	for url := range jobs {
		fmt.Printf("[Worker %d] Starting: %s\n", id, url)

		err := downloadAudio(url, dest, true)

		if err != nil {
			fmt.Printf("[Worker %d] Failed %s: %v\n", id, url, err)
		} else {
			fmt.Printf("[Worker %d] Finished: %s\n", id, url)
		}

	}
}

func main() {
	quiet := flag.Bool("q", false, "Run in quiet mode without progress bar")
	concurrent := flag.Bool("c", false, "Run in concurrent mode with multiple workers")
	flag.Parse()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		return
	}

	dest := filepath.Join(homeDir, "Music", "Download")
	if err := os.MkdirAll(dest, 0755); err != nil {
		fmt.Println("failed to create dir:", err)
		return
	}

	file, err := os.Open("urls.txt")
	if err != nil {
		fmt.Println("failed to open urls.txt:", err)
		return
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("error reading file:", err)
		return
	}

	if len(lines) == 0 {
		fmt.Println("No URLs found in file")
		return
	}

	if *concurrent {
		jobs := make(chan string, len(lines))
		var wg sync.WaitGroup

		numWorkers := 3
		if len(lines) < 3 {
			numWorkers = len(lines)
		}
		fmt.Printf("\n Starting %d workers to process %d links...\n\n", numWorkers, len(lines))

		for i := 1; i <= numWorkers; i++ {
			wg.Add(1)
			go worker(i, jobs, dest, &wg)
		}

		for _, url := range lines {
			jobs <- url
		}
		close(jobs)
		wg.Wait()
		fmt.Println("\n all concurrent downloads over")
	} else {
		fmt.Printf("\nStarting sequential download of %d links...\n\n", len(lines))
		for i, url := range lines {
			fmt.Printf("Downloading %d/%d: %s\n", i+1, len(lines), url)

			if err := downloadAudio(url, dest, *quiet); err != nil {
				fmt.Println(err)
				continue
			}
		}
		fmt.Println("\n all downloads over")

	}
}
