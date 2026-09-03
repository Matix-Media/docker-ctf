#!/usr/bin/env bash
# Raeumt alles auf, was das CTF angelegt hat. Danach faengst du bei Level 1 neu an.
# Cleans up everything the CTF created. Afterwards you restart at level 1.
#
#   bash scripts/ctf-reset.sh

set -u

echo "Raeume auf / cleaning up ..."

docker stop ctf-main data-provider-svc          >/dev/null 2>&1 || true
docker rm   ctf-main data-provider-svc          >/dev/null 2>&1 || true
docker network rm ctf-net                       >/dev/null 2>&1 || true

if [ -d ./secrets ]; then
  rm -rf ./secrets
  echo "  - Ordner ./secrets entfernt / folder ./secrets removed"
fi

echo "  - Container ctf-main und data-provider-svc entfernt / containers removed"
echo "  - Netzwerk ctf-net entfernt / network removed"
echo ""
echo "Fertig. Loesche noch die Cookies fuer localhost:8989 (oder nutze ein"
echo "privates Fenster), damit auch der Fortschrittsbalken zurueckgesetzt ist."
echo ""
echo "Done. Also clear the cookies for localhost:8989 (or use a private window)"
echo "so the progress bar resets too."
echo ""
echo "Neu starten mit / start again with:"
echo "  docker run -d --name ctf-main matixmedia/docker-ctf:latest"
