package nockdir

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"nock/httputils"
	"nock/ioutils"
)

func (d *NockDir) worker(jobs <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for path := range jobs {
		response, err := d.client.Get(path)
		if err != nil {
			fmt.Println("Error making request:", err)
			continue
		}
		d := GetResponseData(response)
		if httputils.IsForbidden(d.StatusCode) {
			continue
		}

		fmt.Printf("%-40s %sStatus: %3d%s [Size: %5d]\n",
			d.Name, utils.StatusColor(d.StatusCode), d.StatusCode, Reset, d.ContentLength)
	}
}

func (d *NockDir) Run(version string) {
	d.Parse(version)
	d.options.DisplayBanner()

	f, err := os.Open(d.options.Wordlist)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer f.Close()

	jobs := make(chan string, d.options.Threads)
	var wg sync.WaitGroup

	// rate := time.Second / 5
	// rateLimiter := time.NewTicker(rate)

	if len(d.options.BaseURL) > 0 && d.options.BaseURL[len(d.options.BaseURL)-1] != '/' {
		d.options.BaseURL += "/"
	}
	for i := 0; i < d.options.Threads; i++ {
		wg.Add(1)
		go d.worker(jobs, &wg)
	}
	go func() {
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			path := scanner.Text()
			if path == "" {
				continue
			}
			if strings.HasPrefix(path, "\uFEFF") {
				path = strings.TrimPrefix(path, "\uFEFF")
			}
			if strings.HasPrefix(path, "#") || strings.HasPrefix(path, "//") || strings.HasPrefix(path, ";") {
				continue
			}
			url := d.options.BaseURL + path
			jobs <- url
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
		}
		close(jobs)
	}()
	wg.Wait()
}
