# === Build Stage ===
# Kompiliert die Go-Anwendung / compiles the Go application
# --platform=$BUILDPLATFORM: der Build laeuft immer nativ auf der Maschine des
# Bauenden. Go cross-kompiliert selbst, wir brauchen also keine QEMU-Emulation.
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

# Von buildx gesetzt, wenn --platform benutzt wird.
ARG TARGETOS=linux
ARG TARGETARCH

WORKDIR /src

# Nur die Modul-Dateien zuerst, damit der Docker-Cache greift
COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o /go/bin/app .

# === Final Stage ===
# Alpine statt scratch: die Teilnehmer sollen mit 'docker exec -it ctf-main sh'
# in den Container schauen koennen - genau so, wie es die Praesentation zeigt.
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /go/bin/app .

# Metadaten fuer Level 2: dieses Label muessen die Teilnehmer finden.
LABEL "ctf.data-provider.host"="data-provider-svc"
LABEL "org.opencontainers.image.title"="Docker CTF"

EXPOSE 8989

# Wird in Level 5 vom Benutzer mit '--show-flag' ueberschrieben.
CMD ["/app/app"]
