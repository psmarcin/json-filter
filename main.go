package main

import (
	"github.com/psmarcin/youtubeGoesPodcast/config"
	"github.com/psmarcin/youtubeGoesPodcast/server"
)

func init() {
	config.Init()
}

func main() {
	server.Start()
}
