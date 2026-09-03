#!/usr/bin/env bash
# Baut beide Images, taggt sie und pusht sie nach Docker Hub.
# Builds both images, tags and pushes them to Docker Hub.
#
#   bash scripts/build-and-push.sh            # baut + pusht :2025 und :latest
#   PUSH=0 bash scripts/build-and-push.sh     # nur bauen, nicht pushen
#   REGISTRY=meinname VERSION=2026 bash scripts/build-and-push.sh

set -euo pipefail

REGISTRY="${REGISTRY:-matixmedia}"
VERSION="${VERSION:-2025}"
PUSH="${PUSH:-1}"

cd "$(dirname "$0")/.."

MAIN="$REGISTRY/docker-ctf-main"
PROV="$REGISTRY/docker-ctf-data-provider"

echo "==> Baue $MAIN:$VERSION"
docker build -t "$MAIN:$VERSION" -t "$MAIN:latest" ./main-app

echo "==> Baue $PROV:$VERSION"
docker build -t "$PROV:$VERSION" -t "$PROV:latest" ./data-provider-app

if [ "$PUSH" != "1" ]; then
  echo "==> PUSH=0, ueberspringe den Push."
  exit 0
fi

echo "==> Pushe nach Docker Hub (docker login noetig)"
docker push "$MAIN:$VERSION"
docker push "$MAIN:latest"
docker push "$PROV:$VERSION"
docker push "$PROV:latest"

echo "==> Fertig."
echo "    $MAIN:$VERSION  /  :latest"
echo "    $PROV:$VERSION  /  :latest"
