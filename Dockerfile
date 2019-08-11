FROM golang:1.12-alpine
ADD . /go/src/topicsync
RUN go build -ldflags "-linkmode external -extldflags -static" -a main.go

FROM scratch
COPY --from=0 /go/src/topicsync /main
CMD ["/main"]