FROM golang:1.11 as base
WORKDIR /go/src/github.com/psmarcin/youtubeGoesPodcast/
COPY . .
RUN go get
RUN make build


FROM gcr.io/distroless/base
COPY --from=base /go/src/github.com/psmarcin/youtubeGoesPodcast/main /
COPY --from=base /go/src/github.com/psmarcin/youtubeGoesPodcast/assets /assets
COPY --from=base /go/src/github.com/psmarcin/youtubeGoesPodcast/templates /templates
EXPOSE 8080
CMD ["/main"]
