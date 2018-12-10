### YoutubeGoesPodcast

> It's simple service that provide feed from youtube channel (in future I'm going to support playlist and users). 


[![CircleCI](https://circleci.com/gh/psmarcin/youtubeGoesPodcast.svg?style=svg)](https://circleci.com/gh/psmarcin/youtubeGoesPodcast) 

Produced feed should help you follow new "videos" in your podcast app. 

There are so many great videos on youtube that works perfect without video, and you only want to listen audio part while doing other things. 

### TODO

* [ ] Improve web ui
* [ ] Add unit tests
* [ ] Endpoint with list for chennals 

### Development

Please keep in mind that this is my side project. I want to learn Go. 


#### Setup 
There are no additional dependency except `go`. 

First please put env variable that you can find in `now.js` into `.env` file and then just run `go run main.go` and you are good to go. 

#### Docker

1. `docker-compose up`
2. `GET localhost:8080` 
