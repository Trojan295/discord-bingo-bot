FROM golang:alpine AS builder

ADD . /bingo-bot/
RUN apk --no-cache add ca-certificates && \
    cd /bingo-bot/ && \
    CGO_ENABLED=0 go build -ldflags="-w -s" -o bingo cmd/bingo/bingo.go


FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /bingo-bot/bingo /opt/bingo

ENTRYPOINT ["/opt/bingo"]

