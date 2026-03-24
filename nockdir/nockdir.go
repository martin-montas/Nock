package nockdir

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"bufio"
	"nock/httputils"
)

const Reset = "\033[0m"
const queueSize = 100
const scannerMaxSize = 1024 * 1024 * 10 // 10 MB

type NockDir struct {
	options *OptionsDir
	client  *httputils.HTTPClient
}

const defaultList = "/usr/share/wordlists/dirb/common.txt"

func (d *NockDir) Parse(version string) {
	if len(os.Args) < 2 {
		fmt.Println("Expected 'dir' , 'crawl' or 'version' subcommand")
		return
	}

	if defaultList == os.Args[1] {
		fmt.Println("Usage: dir -u <url> -w <wordlist> -t <threads>")
		os.Exit(1)
	}

	dirCmd := flag.NewFlagSet("dir", flag.ExitOnError)
	u := dirCmd.String("u", "", "Target URL")
	w := dirCmd.String("w", defaultList, "Wordlist path")
	t := dirCmd.Int("t", 10, "Number of threads")

	if err := dirCmd.Parse(os.Args[2:]); err != nil {
		log.Fatalf("failed to parse dir command: %v", err)
	}

	if *u == "" || *w == "" {
		fmt.Println("Usage: dir -u <url> -w <wordlist> -t <threads>")
		os.Exit(1)
	}

	o := &OptionsDir{
		Wordlist: *w,
		BaseURL:  *u,
		Threads:  *t,
		Version:  version,
	}
	d.options = o
	d.client = httputils.NewHTTPClient()
}

func GetResponseData(r *http.Response) httputils.ResponseData {
	return httputils.ResponseData{
		Name:          r.Request.URL.String(),
		StatusCode:    r.StatusCode,
		ContentLength: r.ContentLength,
	}
}
func streamWords(path string) (<-chan string, <-chan error, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	out := make(chan string, queueSize)
	errc := make(chan error, 1)

	go func() {
		defer close(out)
		defer f.Close()
		scanner := bufio.NewScanner(f)
		// bump buffer size to handle long lines
		buf := make([]byte, 1024*64)
		scanner.Buffer(buf, scannerMaxSize)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			// skip full-line comments
			if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") || strings.HasPrefix(line, ";") {
				continue
			}
			// optional: strip inline comments only if separated by whitespace
			if idx := strings.Index(line, " #"); idx != -1 {
				line = strings.TrimSpace(line[:idx])
			} else if idx := strings.Index(line, "\t#"); idx != -1 {
				line = strings.TrimSpace(line[:idx])
			}
			if line == "" {
				continue
			}

			// will block if channel is full -> backpressure to reader
			out <- line
		}
		if err := scanner.Err(); err != nil {
			// if file too long for buffer, scanner.Err() will tell us
			errc <- err
		}
		close(errc)
	}()

	return out, errc, nil
}
