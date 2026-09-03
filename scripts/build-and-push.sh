#!/usr/bin/env bash
# Baut beide Images und veroeffentlicht sie auf Docker Hub.
# Builds both images and publishes them to Docker Hub.
#
#   bash scripts/build-and-push.sh          # multi-arch bauen + pushen (:2025 und :latest)
#   PUSH=0 bash scripts/build-and-push.sh   # nur lokal bauen, zum Testen
#   REGISTRY=name VERSION=2026 bash scripts/build-and-push.sh
#
# WICHTIG - Multi-Arch:
# Die Teilnehmer-VMs sind in der Regel x86_64 (linux/amd64). Wer auf einem
# Apple-Silicon-Mac baut, erzeugt mit einem normalen 'docker build' aber
# NUR linux/arm64 - die Images liessen sich auf den VMs dann nicht starten.
# Deshalb wird zum Pushen immer buildx mit beiden Plattformen benutzt.

set -euo pipefail

REGISTRY="${REGISTRY:-matixmedia}"
VERSION="${VERSION:-2025}"
PUSH="${PUSH:-1}"
PLATFORMS="${PLATFORMS:-linux/amd64,linux/arm64}"

cd "$(dirname "$0")/.."

MAIN="$REGISTRY/docker-ctf"
PROV="$REGISTRY/docker-ctf-data-provider"

if [ "$PUSH" != "1" ]; then
  echo "==> PUSH=0: baue nur lokal fuer diese Maschine (kein Multi-Arch)"
  docker build -t "$MAIN:$VERSION" -t "$MAIN:latest" ./main-app
  docker build -t "$PROV:$VERSION" -t "$PROV:latest" ./data-provider-app
  echo "==> Fertig. Lokale Images:"
  docker images --format '    {{.Repository}}:{{.Tag}}  {{.ID}}' | grep "$REGISTRY/docker-ctf" || true
  exit 0
fi

# --- Multi-Arch-Builder sicherstellen -------------------------------------
# Der Standard-Builder ('docker' driver) kann kein Multi-Arch. Wir brauchen
# einen Builder mit 'docker-container' driver.
BUILDER="ctf-multiarch"
if ! docker buildx inspect "$BUILDER" >/dev/null 2>&1; then
  echo "==> Lege Buildx-Builder '$BUILDER' an"
  docker buildx create --name "$BUILDER" --driver docker-container >/dev/null
fi
docker buildx inspect --builder "$BUILDER" --bootstrap >/dev/null

echo "==> Angemeldet bei Docker Hub? (docker login noetig)"
docker login

echo "==> Baue und pushe $MAIN  ($PLATFORMS)"
docker buildx build --builder "$BUILDER" \
  --platform "$PLATFORMS" \
  -t "$MAIN:$VERSION" -t "$MAIN:latest" \
  --push ./main-app

echo "==> Baue und pushe $PROV  ($PLATFORMS)"
docker buildx build --builder "$BUILDER" \
  --platform "$PLATFORMS" \
  -t "$PROV:$VERSION" -t "$PROV:latest" \
  --push ./data-provider-app

echo ""
echo "==> Fertig. Veroeffentlichte Plattformen pruefen:"
for img in "$MAIN:latest" "$PROV:latest"; do
  echo "    $img"
  docker manifest inspect "$img" 2>/dev/null \
    | grep -o '"architecture": "[a-z0-9]*"' | sort -u | sed 's/^/      /' || true
done
echo ""
echo "    Danach auf einer x86-Maschine gegenpruefen:"
echo "      docker pull $MAIN:latest && docker run --rm $MAIN:latest --help"
