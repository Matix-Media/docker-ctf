# === Build Stage ===
# --platform=$BUILDPLATFORM: der Build laeuft immer nativ auf der Maschine des
# Bauenden. Go cross-kompiliert selbst, wir brauchen also keine QEMU-Emulation.
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

# Von buildx gesetzt, wenn --platform benutzt wird.
ARG TARGETOS=linux
ARG TARGETARCH

WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /go/bin/dataprovider .

# === Final Stage ===
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /go/bin/dataprovider .

LABEL "org.opencontainers.image.title"="Docker CTF - Data Provider"

EXPOSE 9090

CMD ["/app/dataprovider"]
