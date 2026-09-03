# 🐳 Docker Spickzettel / Cheat Sheet

Eine Seite, alles drauf. Deutsch **und** Englisch.
One page, everything on it. German **and** English.

---

## Images

| Befehl / Command | Deutsch | English |
|---|---|---|
| `docker pull nginx` | Image herunterladen | download an image |
| `docker images` | lokale Images auflisten | list local images |
| `docker build -t name .` | Image aus Dockerfile bauen | build an image from a Dockerfile |

## Container starten / Starting containers

```bash
docker run -d -p 8080:80 --name mein-server nginx
```

| Option | Deutsch | English |
|---|---|---|
| `-d` | läuft im Hintergrund (detached) | runs in the background |
| `-p 8080:80` | Port 8080 bei dir → Port 80 im Container | host port 8080 → container port 80 |
| `--name x` | gibt dem Container einen Namen | gives the container a name |
| `-e VAR=wert` | setzt eine Umgebungsvariable | sets an environment variable |
| `-v /pfad:/pfad` | verbindet einen Ordner mit dem Container | mounts a folder into the container |
| `--network netz` | startet den Container in einem Netzwerk | starts the container on a network |
| `--rm` | löscht den Container automatisch beim Stoppen | removes the container automatically on stop |

## Container anschauen / Looking at containers

| Befehl / Command | Deutsch | English |
|---|---|---|
| `docker ps` | laufende Container | running containers |
| `docker ps -a` | **alle** Container, auch gestoppte | **all** containers, including stopped |
| `docker logs NAME` | Ausgabe eines Containers | output of a container |
| `docker logs -f NAME` | Ausgabe live mitlesen | follow the output live |
| `docker port NAME` | welche Ports sind veröffentlicht? | which ports are published? |

## Container aufräumen / Cleaning up

| Befehl / Command | Deutsch | English |
|---|---|---|
| `docker stop NAME` | Container anhalten | stop a container |
| `docker rm NAME` | gestoppten Container löschen | delete a stopped container |
| `docker rm -f NAME` | stoppen **und** löschen in einem | stop **and** delete in one go |

> 🚨 **Merke / Remember:** Einen Containernamen gibt es nur einmal.
> Vor dem Neustart mit gleichem Namen: `docker stop NAME && docker rm NAME`
> A container name exists only once. Before restarting with the same name, remove it.

## In den Container schauen / Looking inside

| Befehl / Command | Deutsch | English |
|---|---|---|
| `docker inspect NAME` | alle Metadaten als JSON (sehr lang!) | all metadata as JSON (very long!) |
| `docker inspect --format '{{.Config.ExposedPorts}}' NAME` | nur die Ports | just the ports |
| `docker inspect --format '{{.Config.Labels}}' NAME` | nur die Labels | just the labels |
| `docker exec -it NAME sh` | Shell im Container öffnen | open a shell in the container |
| `docker exec -it NAME BEFEHL` | Befehl im Container ausführen | run a command in the container |

`-i` = interactive (Tastatur kommt an / keyboard reaches it)
`-t` = TTY (echtes Terminal / real terminal)

## Netzwerke / Networks

| Befehl / Command | Deutsch | English |
|---|---|---|
| `docker network create netz` | Netzwerk anlegen | create a network |
| `docker network ls` | Netzwerke auflisten | list networks |
| `docker run --network netz ...` | Container ins Netzwerk hängen | start a container on the network |

> Im **eigenen** Netzwerk wird der `--name` zum Hostnamen.
> Im **Standard**-Netzwerk funktionieren Namen **nicht**.
> On your **own** network the `--name` becomes the hostname.
> On the **default** network, names do **not** work.

## Volumes / Bind Mounts

```bash
docker run -v $(pwd)/daten:/app/daten meinimage
```

- Links: Ordner **auf deinem Rechner** — Docker braucht einen **vollständigen** Pfad.
- Rechts: Ordner **im Container**.
- `$(pwd)` = der Ordner, in dem du gerade stehst.

Left: folder **on your machine** — Docker needs a **full** path.
Right: folder **in the container**. `$(pwd)` = the folder you are currently in.

## Docker Compose

| Befehl / Command | Deutsch | English |
|---|---|---|
| `docker compose up -d` | alle Dienste starten | start all services |
| `docker compose down` | alles stoppen und entfernen | stop and remove everything |
| `docker compose ps` | Status aller Dienste | status of all services |

---

## Wenn etwas nicht geht / When something breaks

| Fehlermeldung / Error | Ursache / Cause | Lösung / Fix |
|---|---|---|
| `container name is already in use` | alter Container existiert noch / old container still exists | `docker rm -f NAME` |
| `port is already allocated` | Port schon belegt / port already taken | anderen Port nehmen oder `docker ps` prüfen |
| `pull access denied` | Image-Name falsch / wrong image name | Schreibweise prüfen / check the spelling |
| `no such host` | Name nicht auflösbar / name not resolvable | gleiches Netzwerk? / same network? |
| `connection refused` | niemand lauscht / nobody listening | läuft der Container? / is the container running? |
| Terminal hängt / terminal hangs | `-d` vergessen / forgot `-d` | `Strg+C`, dann mit `-d` neu starten |
