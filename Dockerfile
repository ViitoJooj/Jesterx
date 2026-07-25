FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
RUN CGO_ENABLED=0 go build -o /out/daemon ./cmd/api

FROM alpine:3.20
COPY --from=build /out/daemon /usr/local/bin/daemon
COPY migrations /migrations
EXPOSE 8080
ENTRYPOINT ["daemon"]
