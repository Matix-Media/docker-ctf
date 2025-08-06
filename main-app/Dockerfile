# === Build Stage ===
# Kompiliert die Go-Anwendung
FROM golang:1.24-alpine AS builder

WORKDIR /src

# Kopiere nur die Modul-Dateien, um den Docker-Cache zu nutzen
COPY go.mod go.sum* ./
RUN go mod download

# Kopiere den Rest des Codes und baue die Anwendung
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/app .

# === Final Stage ===
# Erstellt das finale, minimale Image
FROM scratch

WORKDIR /app

# Kopiere nur die kompilierte Anwendung aus dem Build-Stage
COPY --from=builder /go/bin/app .

# Metadaten für die CTF-Schritte
LABEL "ctf.data-provider.host"="data-provider-svc"
EXPOSE 8989

# Startbefehl für den Container
# Dieser wird im letzten Schritt vom Benutzer mit '--show-flag' überschrieben.
CMD ["/app/app"]
