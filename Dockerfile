FROM golang:1.12-alpine as builder

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create unpriv user to run bot
RUN adduser -D -g '' appuser

WORKDIR $GOPATH/src/topicsync
COPY . .

RUN go get -d -v

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/topicsync .

FROM scratch

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/topicsync /go/bin/topicsync
# Copy over the root certs from ca-certificates
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
USER appuser

CMD ["/go/bin/topicsync"]
