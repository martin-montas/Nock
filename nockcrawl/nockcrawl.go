package nockcrawl

import (
	"flag"
	"fmt"
	"log"
	"os"

	"nock/httputils"
)

const Reset = "\033[0m"

type NockCrawl struct {
	opt    *OptionsCrawl
	client *httputils.HTTPClient
	href   *[]Href
}

// parses the command line for the current
// type of module
func (c *NockCrawl) Parse(version string) {
	if len(os.Args) < 2 {
		fmt.Println("Expected 'dir' , 'crawl' or 'version' subcommand")
		return
	}
	dirCmd := flag.NewFlagSet("dir", flag.ExitOnError)
	u := dirCmd.String("u", "", "Target URL")
	t := dirCmd.Int("t", 10, "Number of threads")
	if err := dirCmd.Parse(os.Args[2:]); err != nil {
		log.Fatalf("failed to parse dir command: %v", err)
	}

	if *u == "" {
		fmt.Println("Usage: dir -u <url> -w <wordlist> -t <threads>")
		return
	}
	o := &OptionsCrawl{
		BaseURL: *u,
		Threads: *t,
		Version: version,
	}
	c.opt = o
	c.client = httputils.NewHTTPClient()
}
