
FROM open-match-base-build as builder
ADD . /go/src/open-match-example
COPY demo/static /app/static
WORKDIR /go/src/open-match-example

ENV GO111MODULE on

ENTRYPOINT ["go", "run", "main.go"]
