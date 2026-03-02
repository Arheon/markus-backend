FROM golang:1.25 as base

COPY ./ /opt/app
WORKDIR /opt/app

RUN go install github.com/go-delve/delve/cmd/dlv@v1.21
RUN go mod download

ENTRYPOINT ["sh", "-c", "while true; do sleep 1; done"]
