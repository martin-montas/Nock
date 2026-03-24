package nockcrawl

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
	"golang.org/x/net/html"
	"nock/httputils"
)

var tags = []string{"a", "link", "base", "area"}

func ExtractLinksFromNode(n *html.Node, baseURL url.URL) []Href {
	links := []Href{}
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			parsed, err := url.Parse(attr.Val)
			if err != nil {
				fmt.Println("Bad URL:", err)
				continue
			}
			resolved := baseURL.ResolveReference(parsed)
			response, err := httputils.FetchResponse(resolved.String())
			response.Body.Close()
			if err != nil {
				continue
			}
			if SameDomain(baseURL.String(), resolved.String()) && response.StatusCode == 200 {
				links = append(links, Href{
					StatusCode: response.StatusCode,
					Path:       resolved.String(),
					Alive:      true,
				})
			} else {
				continue
			}
		}
	}
	return links
}

func extractRecursive(doc *html.Node, baseURL url.URL, links []Href) []Href {
	if doc.Type == html.ElementNode && contains(tags, doc.Data) {
		if whichSection(doc) != "head" {
			newLinks := ExtractLinksFromNode(doc, baseURL)
			for _, link := range newLinks {
				if link.Alive {
					links = append(links, link)
				}
			}
		}
	}
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		links = extractRecursive(c, baseURL, links)
	}

	return links
}

func whichSection(n *html.Node) string {
	// Determine if node is in <head> or <body>
	for p := n.Parent; p != nil; p = p.Parent {
		if p.Type == html.ElementNode {
			if p.Data == "head" {
				return "head"
			}
			if p.Data == "body" {
				return "body"
			}
		}
	}
	return "unknown"
}

// SameDomain checks if urlA and urlB are in the same domain or subdomain relationship
func SameDomain(urlA, urlB string) bool {
	parsedA, err1 := url.Parse(urlA)
	parsedB, err2 := url.Parse(urlB)
	if err1 != nil || err2 != nil {
		return false
	}
	hostA := strings.ToLower(parsedA.Hostname())
	hostB := strings.ToLower(parsedB.Hostname())

	return hostA == hostB || strings.HasSuffix(hostA, "."+hostB) || strings.HasSuffix(hostB, "."+hostA)
}

func ExtractLinks(doc *html.Node, baseURL url.URL) []Href {
	var links []Href
	return extractRecursive(doc, baseURL, links)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func worker(wg *sync.WaitGroup, parsedURL *url.URL, set mapset.Set[string], l []Href) {
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
