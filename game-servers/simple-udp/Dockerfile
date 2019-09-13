FROM golang:1.11.5 as builder
WORKDIR /go/src/simple-udp
RUN go get agones.dev/agones/pkg/sdk
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server .

# final image
FROM alpine:3.9

RUN adduser -D server
COPY --from=builder /go/src/simple-udp/server /home/server/server
RUN chown -R server /home/server && \
    chmod o+x /home/server/server

USER server
ENTRYPOINT ["/home/server/server"]