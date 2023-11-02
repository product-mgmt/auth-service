FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache make
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN make build

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/target/app-linux .
COPY --from=builder /app/.env .
EXPOSE 9001

CMD [ "/app-linux" ]
