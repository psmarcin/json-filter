FROM golang:1.11-alpine3.8 as build
WORKDIR /app
COPY . .
RUN ./make


FROM alpine:3.7
WORKDIR /app
COPY --from=build /app /app/
EXPOSE 8080
CMD ["./run"]
