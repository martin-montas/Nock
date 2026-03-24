package cmd

import (
	"nock/nockcrawl"
	// "nock/nockfuzz"
	"nock/nockdir"
)

var Version = "1.0.0"

type Parser interface {
	Run(version string)
}

// modules should be call here
var registry = map[string]Parser{
	"dir": &nockdir.NockDir{},
	"crawl": &nockcrawl.NockCrawl{},
	// "fuzz":  &nockfuzz.NockFuzz{},
}

func Parse(c Parser) {
	c.Run(Version)
}
