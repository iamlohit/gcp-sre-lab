FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go mod init sre-app && go mod tidy && go build -o app .

FROM gcr.io/distroless/base
COPY --from=builder /app/app /
CMD ["/app"]