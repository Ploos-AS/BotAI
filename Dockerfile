FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go test ./... && CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/botai ./cmd/botai

FROM alpine:3.22
RUN addgroup -S botai && adduser -S -G botai botai
COPY --from=build /out/botai /usr/local/bin/botai
USER botai
EXPOSE 8090
ENV BOTAI_LISTEN=0.0.0.0:8090
ENTRYPOINT ["/usr/local/bin/botai"]
