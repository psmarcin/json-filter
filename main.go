package main

import (
	"github.com/psmarcin/youtubeGoesPodcast/pkg/config"
	"github.com/psmarcin/youtubeGoesPodcast/pkg/server"
)

func init() {
	config.Init()
}

func main() {
	server.Start()
}
