# Batch Music Downloader

A fast, concurrent command-line interface (CLI) application written in Go that reads a local list of URLs and batch-downloads them as MP3s. Powered by Go's native concurrency and [yt-dlp](https://github.com/yt-dlp/yt-dlp).

## Features

Reads links from a `urls.txt` file and safely downloads them into `~/Music/Download`, featuring:
- **Custom UI:** A live, animated terminal progress bar.
- **Execution Modes:** Run sequentially for a clean UI, or concurrently for maximum speed.
- **Resiliency:** Bypasses broken links without crashing the entire batch process.
- **Smart Syncing:** Automatically skips files that have already been downloaded.

## Usage 

You can run the program directly using Go. Ensure you have a `urls.txt` file in the same directory with one URL per line.

**Standard Mode (Visual Progress Bar):**

```bash
go run main.go
```

**Quiet Mode (No Progress Bar):**

```Bash
go run main.go -q
```

**Concurrent Mode (Spawns 3 background workers):**

```Bash
go run main.go -c
```

Example Output (Concurrent Mode):

```Plaintext
🚀 Starting 3 workers to process 11 links...

[Worker 3] Starting: https://youtu.be/dvgZkm1xWPE?si=eogtrb8hyq8-7JB9
[Worker 1] Starting: https://youtu.be/nLnp0tpZ0ok?si=1g3bV6BqgH2q9MvT
[Worker 2] Starting: https://youtu.be/AGsn2ycFRqI?si=3WvjCBGY09KZwABx
[Worker 1] Finished: https://youtu.be/nLnp0tpZ0ok?si=1g3bV6BqgH2q9MvT
[Worker 1] Starting: https://youtu.be/m7Bc3pLyij0?si=1kK4-9jPJ0MlEgs_
[Worker 2] Finished: https://youtu.be/AGsn2ycFRqI?si=3WvjCBGY09KZwABx
[Worker 2] Starting: https://youtu.be/kPtn26x8TZM?si=VMLyuxEGR6js5MQ5
[Worker 3] Finished: https://youtu.be/dvgZkm1xWPE?si=eogtrb8hyq8-7JB9
[Worker 3] Starting: https://youtu.be/NobzfIebbrE?si=RmLb3N1ddwEpwQfX
...
🎉 All concurrent downloads complete!
```

## Build it yourself

To compile the app into a standalone executable that you can run from anywhere on your system:

```Bash
go build -o bin/music-downloader
./music-downloader -c
```

## Under the Hood

- Subprocess Control: Uses Go's os/exec package to safely trigger the yt-dlp binary, explicitly managing its working directory via cmd.Dir to handle cross-platform file paths cleanly.

- Standard Output Interception: Instead of blocking execution, the app uses cmd.StdoutPipe() and a bufio.Scanner to intercept the stdout stream of the subprocess in real-time.

- Regex UI Parsing: Uses Go's regexp package (compiled globally at startup for performance) to parse percentage values out of the raw yt-dlp text stream, utilizing carriage returns (\r) to draw a custom, single-line animated progress bar.

- Concurrency (Worker Pool): Implements a highly efficient goroutine worker pool using sync.WaitGroup and channels. This decouples the file reading (job generation) from the actual downloading (job execution).

```Go
// Compiling regex once at global scope to save CPU cycles during the scanner loop
var progressRegex = regexp.MustCompile(`\[download\]\s+([0-9.]+)%`)
```