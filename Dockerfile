FROM golang:1.11-alpine3.8 as build
WORKDIR /go/src/github.com/psmarcin/json-filter/
COPY . .
RUN ./make


FROM alpine:3.7
WORKDIR /app
COPY --from=build /go/src/github.com/psmarcin/json-filter /app/
EXPOSE 8080
CMD ["./run"]
