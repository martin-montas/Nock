package nockdir

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"nock/httputils"
)

const Reset = "\033[0m"

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
