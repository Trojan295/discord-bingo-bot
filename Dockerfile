FROM golang:alpine AS builder

ADD . /bingo-bot/
RUN cd /bingo-bot/ && \
    CGO_ENABLED=0 go build -ldflags="-w -s" -o bingo cmd/bingo/bingo.go


FROM scratch
COPY --from=builder /bingo-bot/bingo /opt/bingo

ENTRYPOINT ["/opt/bingo"]

