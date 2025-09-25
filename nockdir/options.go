package nockdir

import (
	"fmt"
)

type OptionsDir struct {
	Wordlist string
	BaseURL  string
	Threads  int
	Version  string
}

// displays the initial information
func (o *OptionsDir) DisplayBanner() {
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
`, o.Version, o.BaseURL, o.Wordlist, o.Threads)
}
