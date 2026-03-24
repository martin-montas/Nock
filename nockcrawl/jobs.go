package nockcrawl

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"

	"nock/httputils"

	mapset "github.com/deckarep/golang-set/v2"
	"golang.org/x/net/html"
)

type Href struct {
	StatusCode int
	Path       string
	Alive      bool
}

func (c *NockCrawl) worker(wg *sync.WaitGroup, parsedURL *url.URL, set mapset.Set[string]) {
	defer wg.Done()

	for i := 0; i < len(l); i++ {
		resp, err := http.Get(l[i].Path)
		if err != nil {
			fmt.Printf("Domain is unreachable: %s\n", l[i].Path)
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("Failed to read body: %s\n", l[i].Path)
			continue
		}
		doc, err := html.Parse(bytes.NewReader(body))
		if err != nil {
			fmt.Println("HTML parse error:", err)
			continue
		}
		links := ExtractLinks(doc, *parsedURL)
		for _, link := range links {
			if set.Contains(link.Path) {
				continue
			}
			if !set.Contains(link.Path) {
				set.Add(link.Path)
				l = append(l, link)
			}
		}
	}
}

func (c *NockCrawl) Run(version string) {
	c.Parse(version)
	c.opt.DisplayBanner()
	baseURLStatus := httputils.CheckStatuscodeFromURL(c.opt.BaseURL)
	if baseURLStatus != 200 {
		fmt.Printf("url is unreachable %s\n", c.opt.BaseURL)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Println("Usage: crawl -u <url> -t <threads>")
		os.Exit(1)
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Println("Usage: crawl -u <url> -t <threads>")
		os.Exit(1)
	}

	if os.Args[1] == "version" || os.Args[1] == "-v" || os.Args[1] == "--version" {
		fmt.Println("version:", version)
		os.Exit(1)
	}
	if c.opt.BaseURL[len(c.opt.BaseURL)-1:] != "/" {
		c.opt.BaseURL += "/"
	}

	parsedURL, err := url.Parse((c.opt.BaseURL))
	if err != nil {
		panic(err)
	}

	set := mapset.NewSet[string]()

	Links := []Href{
		{
			StatusCode: 200,
			Path:       c.opt.BaseURL,
			Alive:      true,
		},
	}

	var wg sync.WaitGroup
	wg.Add(1)
	jobs := make(chan *url.URL)
	go worker(&wg, parsedURL, set, Links)

	jobs <- parsedURL
	wg.Wait()
	it := set.Iterator()

	for value := range it.C {
		fmt.Printf("Discovered: %s\n", value)
	}
}
