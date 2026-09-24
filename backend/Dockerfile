# ---- Build stage ----
FROM golang:1.22-alpine AS build
WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

# ---- Run stage ----
FROM alpine:3.20
RUN adduser -D -g '' appuser
COPY --from=build /out/server /usr/local/bin/server

USER appuser
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/server"]
