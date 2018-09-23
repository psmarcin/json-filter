FROM golang:1.11-alpine3.8 as build
WORKDIR /go/src/github.com/psmarcin/youtubeGoesPodcast/
COPY . .
RUN ./make.sh


FROM alpine:3.7
WORKDIR /app
COPY --from=build /go/src/github.com/psmarcin/youtubeGoesPodcast /app/
EXPOSE 8080
CMD ["./run.sh"]
