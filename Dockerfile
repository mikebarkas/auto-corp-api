FROM golang:1.21 as build
WORKDIR /
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build main.go

FROM alpine:3.18
WORKDIR /
COPY --from=build /main /main
EXPOSE 8080
ENTRYPOINT ["/main"]
