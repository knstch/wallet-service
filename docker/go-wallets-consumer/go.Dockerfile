FROM golang:1.24 AS base

FROM base AS builder

WORKDIR /build
COPY . ./
RUN go build ./cmd/consumer

FROM base AS final

ARG PORT

WORKDIR /app
COPY --from=builder /build/consumer /build/.env ./
COPY --from=builder /build/consumer ./

CMD ["/app/consumer"]