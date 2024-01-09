# syntax=docker/dockerfile:1

FROM golang:1.21.4

WORKDIR /home/my/go/src/github.com/udonetsm/cmngb

COPY ./ ./
RUN go build -o ./cmngb .
CMD ["./cmngb"]