package utils

import "fmt"

const (
	Red   = "\033[31m"
	Green = "\033[32m"
	Blue  = "\033[34m"
	Reset = "\033[0m"
)

func InitialInfoDirForcer(url string, wordlist string, threads int, version string) {
	fmt.Printf(`
===============================================================
Nock %-6s 
by martin montas - @github.com/martin-montas
===============================================================
[+] URL:      		%-21s
[+] Wordlist: 		%-21s
[+] Threads: 		%-21d
===============================================================
                       STARTING                       
===============================================================
`, version, url, wordlist, threads)
}

func InitialInfoCrawler(url string, threads int, version string) {
	fmt.Printf(`
===============================================================
Nock %-6s 
by martin montas - @github.com/martin-montas
===============================================================
[+] URL:      		%-21s
[+] Threads: 		%-21d
===============================================================
                       STARTING                       
===============================================================
`, version, url, threads)
}

func StatusColor(status int) string {
	switch {
	case status >= 100 && status < 200:
		return "\033[36m" // Cyan
	case status >= 200 && status < 300:
		return "\033[32m" // Green
	case status >= 300 && status < 400:
		return "\033[34m" // Blue
	case status >= 400 && status < 500:
		return "\033[33m" // Yellow
	case status >= 500:
		return "\033[31m" // Red
	default:
		return "\033[35m" // Magenta for unknown
	}
}
