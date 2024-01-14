# syntax=docker/dockerfile:1

FROM golang:1.21.4

WORKDIR $GOPATH/src/github.com/udonetsm/cmngb
COPY ./ ./
COPY ./database/cfg.yaml /etc/
RUN go build -o ./cmngb .
CMD ["./cmngb"]