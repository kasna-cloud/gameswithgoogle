FROM open-match-base-build as builder
ADD . /go/src/open-match-function-pool
WORKDIR /go/src/open-match-function-pool

ENV GO111MODULE on

ENTRYPOINT ["go", "run", "main.go"]