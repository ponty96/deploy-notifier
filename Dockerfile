# ------------------------------------------------------------------------------
# Build Container
# ------------------------------------------------------------------------------
FROM golang:1.22.3 as builder

ENV GOPATH=/go
ENV GO111MODULE=on

ADD . /deploy-notifier
WORKDIR /deploy-notifier

RUN go mod download

# RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags netgo -ldflags '-w -extldflags "-static"'
RUN go build -o main .

# ------------------------------------------------------------------------------
# Production Container
# ------------------------------------------------------------------------------
FROM alpine:3.20
COPY --from=builder /deploy-notifier/deploy-notifier /deploy-notifier
RUN apk add ca-certificates

CMD ["./deploy-notifier"]
