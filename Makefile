current_dir := $(shell pwd)

# PROD
build:
	CGO_ENABLED=0 go build -ldflags "-s -w" -o main

run:
	./main

start:
	docker build -t ygp-prod .
	docker run -it -p 8080:8080 ygp-prod
	
# DEV
dev:
	docker build -t ygp -f Dockerfile-dev .
	docker run -it -p 8080:8080 -v ${current_dir}:/go/src/github.com/psmarcin/youtubeGoesPodcast ygp

dev-build:
	docker build -t ygp -f Dockerfile-dev .

dev-run:
	docker run -it -p 8080:8080 -v ${current_dir}:/go/src/github.com/psmarcin/youtubeGoesPodcast ygp

# DEPLOY
deploy:
	now
	now alias podcast.psmarcin.me
	now rm podcast --safe --yes 

deploy-ci:
	now -t ${NOWSHTOKEN}
	now alias podcast.psmarcin.me -t ${NOWSHTOKEN}
