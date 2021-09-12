FROM golang:1.15.6-alpine AS builder
RUN apk add --update alpine-sdk gcc make git vim bash curl openssh autoconf automake gettext libtool  && rm -rf /var/cache/apk/*
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -i -v -o ./cmd/run/pow ./cmd/run

FROM alpine:3.12.3 AS runtime
LABEL vendor="Kirill Govina"
LABEL maintainer="Kirill Govina <radiophysic@gmail.com>"
LABEL description="Server | Proof-of-work concept"
RUN apk add --update libstdc++ libgcc gettext gnu-libiconv && rm -rf /var/cache/apk/*
WORKDIR /app
COPY --from=builder     /app/cmd/run/pow ./
COPY --from=builder     /app/assets      ./assets/
CMD ["./pow -mode=server"]
EXPOSE 7777
