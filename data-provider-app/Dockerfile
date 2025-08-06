# === Build Stage ===
FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/dataprovider .

# === Final Stage ===
FROM scratch

WORKDIR /app
COPY --from=builder /go/bin/dataprovider .

EXPOSE 9090

CMD ["/app/dataprovider"]
