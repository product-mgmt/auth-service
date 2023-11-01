FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o authservice cmd/main.go

FROM gcr.io/distroless/base-debian12
COPY --from=builder /app/authservice .
EXPOSE 9001

CMD [ "/authservice" ]
