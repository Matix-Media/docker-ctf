# 🐳 Docker Capture the Flag

Ein spielerischer Einstieg in Docker: fünf Level, fünf Flags, kein Vorwissen nötig.
A playful introduction to Docker: five levels, five flags, no prior knowledge required.

> ℹ️ **Dieses Repository ist für Trainer und Maintainer.**
> **Spieler bekommen ausschließlich den Image-Namen** — keinen Repo-Link, keine Anleitung.
> Alles, was sie brauchen, steht in den Container-Logs und in der Weboberfläche.
> Hier im Repo liegt die Komplettlösung.
>
> **This repository is for trainers and maintainers.** Players are given *only the image
> name* — no repo link, no instructions. Everything they need is in the container logs and
> the web UI. The full solution lives here.

---

## 🎮 Was die Spieler bekommen / What players get

Genau eine Zeile:

```
matixmedia/docker-ctf
```

Mehr nicht. Der Container erklärt sich ab da selbst:

1. Sie starten ihn und lesen die Logs → dort steht Level 1 komplett erklärt.
2. Ab Level 2 übernimmt die Weboberfläche: Fortschrittsanzeige, verdiente Flags,
   das aktuelle Level und drei aufklappbare Tipps pro Level.
3. Das zweite Image (`data-provider`) entdecken sie **im Spiel** — es steht bewusst
   nirgends vorab.

Die Oberfläche gibt es auf **Deutsch und Englisch** (Umschalter oben rechts), in
**einem** Image — die Befehle sind für alle identisch.

---

## 🧭 Die fünf Level

| Level | Titel | Docker-Konzept |
|---|---|---|
| 1 | Der geheime Port | `docker inspect --format`, `-p` |
| 2 | Der unsichtbare Freund | Image-Labels, `--name` |
| 3 | Das Netzwerk | `docker network create`, `--network`, DNS |
| 4 | Der Tresor | Bind Mount `-v`, `$(pwd)` |
| 5 | Die Flagge | `docker exec -it`, TTY |

Jedes Level vergibt eine eigene Flag. Spieler posten sie in Teams — so sieht der Trainer
live, wer wo steht.

Davor gehört eine **Aufwärmrunde (Level 0)** mit `nginx`, die der Trainer gemeinsam mit
allen vorne durchmacht. Siehe Trainer-Guide.

---

## 🧑‍🏫 Für Trainer / For trainers

| Datei | Zweck |
|---|---|
| **[docs/TRAINER.md](docs/TRAINER.md)** | **Hier anfangen.** Ablauf, Timing, Pre-Flight-Check, Fehlerbilder |
| [docs/SOLUTION.de.md](docs/SOLUTION.de.md) · [en](docs/SOLUTION.en.md) | Komplettlösung mit Erklärungen |
| [docs/SLIDES-ADDENDUM.md](docs/SLIDES-ADDENDUM.md) | Fehlende Folien für die Präsentation |
| [docs/CHEATSHEET.md](docs/CHEATSHEET.md) | Spickzettel, zweisprachig |
| [docs/WORKBOOK.de.md](docs/WORKBOOK.de.md) · [en](docs/WORKBOOK.en.md) | Optionales Arbeitsheft (siehe unten) |

### Zum Workbook

Das CTF ist **ohne** Workbook lösbar — der Container erklärt alles selbst. Das Workbook
ist ein **optionaler Ausdruck** für Teilnehmer, die lieber etwas auf Papier haben, und es
enthält als einziges die **Aufwärmrunde Level 0**. Wenn du es austeilst, teile es als PDF
oder Ausdruck aus, **nicht als Repo-Link** — im Repo steht die Lösung daneben.

---

## 🔧 Entwicklung / Development

```bash
docker compose up --build
```

Baut und startet beides lokal auf <http://localhost:8989>, ohne Docker Hub.

> ⚠️ Compose konfiguriert Netzwerk, Namen und Volume bereits fertig — damit sind genau die
> Aufgaben gelöst, die die Teilnehmer selbst machen sollen. Nur zum Entwickeln benutzen.

```bash
bash scripts/build-and-push.sh
```

Baut, taggt (`:2025` **und** `:latest`) und pusht beide Images. `docker login` nötig.
Nur bauen: `PUSH=0 bash scripts/build-and-push.sh`.
Anderer Account / neues Jahr: `REGISTRY=name VERSION=2026 bash scripts/build-and-push.sh`.

### Images

| Image | Rolle | Spieler erfahren davon … |
|---|---|---|
| `matixmedia/docker-ctf` | Hauptcontainer, Level 1–5 | vorab (das Einzige, was sie bekommen) |
| `matixmedia/docker-ctf-data-provider` | zweiter Container, Level 3–4 | erst im Spiel, in Level 3 |

### Hilfsskripte

```bash
bash scripts/ctf-status.sh   # was läuft, was fehlt, was ist der nächste Schritt
bash scripts/ctf-reset.sh    # Container, Netzwerk und ./secrets aufräumen
```

Für Trainer beim Debuggen gedacht. Spieler brauchen sie nicht — die Weboberfläche nennt
die Aufräum-Befehle selbst.
