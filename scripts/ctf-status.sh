#!/usr/bin/env bash
# Zeigt, was gerade laeuft, und sagt dir, was der naechste Schritt ist.
# Shows what is running and tells you what to do next.
#
#   bash scripts/ctf-status.sh

set -u

running() { docker ps --format '{{.Names}}' | grep -qx "$1"; }
exists()  { docker ps -a --format '{{.Names}}' | grep -qx "$1"; }

echo "=============================================="
echo " Docker CTF — Status"
echo "=============================================="
echo ""

# --- Docker selbst
if ! docker info >/dev/null 2>&1; then
  echo "  Docker laeuft nicht / Docker is not running."
  echo "  -> Starte Docker und versuche es nochmal."
  exit 1
fi

# --- ctf-main
if running ctf-main; then
  PORTS=$(docker port ctf-main 2>/dev/null | tr '\n' ' ')
  echo "  ctf-main            : laeuft / running"
  if [ -n "$PORTS" ]; then
    echo "                        Ports: $PORTS"
  else
    echo "                        !! kein Port veroeffentlicht / no port published"
  fi
  NET=$(docker inspect --format '{{range $k,$v := .NetworkSettings.Networks}}{{$k}} {{end}}' ctf-main)
  echo "                        Netzwerk / network: $NET"
  MNT=$(docker inspect --format '{{range .Mounts}}{{.Destination}} {{end}}' ctf-main)
  echo "                        Mounts: ${MNT:-keine / none}"
elif exists ctf-main; then
  echo "  ctf-main            : existiert, laeuft aber nicht / exists but stopped"
else
  echo "  ctf-main            : nicht vorhanden / not present"
fi

# --- data-provider-svc
if running data-provider-svc; then
  NET2=$(docker inspect --format '{{range $k,$v := .NetworkSettings.Networks}}{{$k}} {{end}}' data-provider-svc)
  echo "  data-provider-svc   : laeuft / running   (Netzwerk: $NET2)"
elif exists data-provider-svc; then
  echo "  data-provider-svc   : existiert, laeuft aber nicht / exists but stopped"
else
  echo "  data-provider-svc   : nicht vorhanden / not present"
fi

# --- Netzwerk
if docker network ls --format '{{.Name}}' | grep -qx ctf-net; then
  echo "  Netzwerk ctf-net    : vorhanden / exists"
else
  echo "  Netzwerk ctf-net    : nicht vorhanden / not present"
fi

# --- secrets
if [ -f ./secrets/password.txt ]; then
  echo "  secrets/password.txt: vorhanden / exists"
else
  echo "  secrets/password.txt: nicht vorhanden / not present"
fi

echo ""
echo "----------------------------------------------"
echo " Naechster Schritt / next step"
echo "----------------------------------------------"

if ! exists ctf-main; then
  echo "  Level 1: Starte den Hauptcontainer und lies die Logs."
  echo "    docker run -d --name ctf-main matixmedia/docker-ctf:latest"
  echo "    docker logs ctf-main"
elif ! running ctf-main; then
  echo "  ctf-main ist gestoppt. Aufraeumen und neu starten:"
  echo "    docker rm ctf-main"
elif ! docker port ctf-main 2>/dev/null | grep -q 8989; then
  echo "  Level 1: Port 8989 ist noch nicht veroeffentlicht."
  echo "    docker stop ctf-main && docker rm ctf-main"
  echo "    docker run -d --name ctf-main -p 8989:8989 matixmedia/docker-ctf:latest"
elif ! running data-provider-svc; then
  echo "  Level 2/3: Der Data-Provider laeuft noch nicht."
  if ! docker network ls --format '{{.Name}}' | grep -qx ctf-net; then
    echo "    docker network create ctf-net"
  fi
  echo "    docker run -d --name data-provider-svc --network ctf-net \\"
  echo "      matixmedia/docker-ctf-data-provider:latest"
elif ! docker inspect --format '{{range $k,$v := .NetworkSettings.Networks}}{{$k}} {{end}}' ctf-main | grep -q ctf-net; then
  echo "  Level 3: ctf-main ist nicht im Netzwerk ctf-net. Beide Container muessen drin sein!"
  echo "    docker stop ctf-main && docker rm ctf-main"
  echo "    docker run -d --name ctf-main --network ctf-net -p 8989:8989 \\"
  echo "      matixmedia/docker-ctf:latest"
elif [ ! -f ./secrets/password.txt ]; then
  echo "  Level 4: Der Ordner /secrets ist noch nicht gemountet."
  echo "    docker stop ctf-main && docker rm ctf-main"
  echo "    docker run -d --name ctf-main --network ctf-net -p 8989:8989 \\"
  echo "      -v \$(pwd)/secrets:/secrets matixmedia/docker-ctf:latest"
else
  echo "  Level 5: Alles steht! Hol dir die Flagge:"
  echo "    docker exec -it ctf-main /app/app --show-flag"
fi
echo ""
echo "  Web: http://localhost:8989"
echo ""
