# 🐳 Docker Capture the Flag

**[Deutsch](#deutsch) · [English](#english)**

Ein spielerischer Einstieg in Docker: fünf Level, fünf Flags, kein Vorwissen nötig.
A playful introduction to Docker: five levels, five flags, no prior knowledge required.

---

<a name="deutsch"></a>
## 🇩🇪 Deutsch

### Was ist das hier?

Du löst nach und nach fünf Level. In jedem Level lernst du **einen** neuen Docker-Befehl
kennen und schaltest damit eine **Flag** frei — einen kurzen Code wie `FLAG{...}`.
Am Ende wartet die große finale Flag.

Du musst **nichts programmieren**. Du tippst Befehle in ein Terminal und schaust dir das
Ergebnis im Browser an.

### Was du brauchst

- Einen Rechner mit **Docker** (in der Bootcamp-VM ist alles schon da).
- Ein **Terminal**.
- Einen **Browser**.
- **Rund 90 Minuten** Zeit.

Prüfe zuerst, ob Docker läuft:

```bash
docker --version
```

Kommt eine Versionsnummer? Perfekt. Kommt ein Fehler, sag kurz Bescheid.

### Was du am Ende kannst

| Level | Das lernst du |
|---|---|
| 0 | Container starten, anschauen, stoppen und aufräumen |
| 1 | In ein Image hineinschauen (`docker inspect`) und Ports freigeben (`-p`) |
| 2 | Metadaten (Labels) lesen und Container benennen (`--name`) |
| 3 | Container miteinander reden lassen (`docker network`) |
| 4 | Dateien zwischen Rechner und Container teilen (`-v`) |
| 5 | Befehle in einem laufenden Container ausführen (`docker exec -it`) |

### Los geht's

Hol dir zuerst die beiden Images:

```bash
docker pull matixmedia/docker-ctf-main:latest
```

```bash
docker pull matixmedia/docker-ctf-data-provider:latest
```

Und dann öffne das Workbook — da steht alles Schritt für Schritt drin:

### ➡️ **[docs/WORKBOOK.de.md](docs/WORKBOOK.de.md)** ⬅️

> 💡 Die Weboberfläche des CTF gibt es auf Deutsch und Englisch. Oben rechts umschalten.

### Hilfe

- **Ich stecke fest** → Im Workbook hat jedes Level drei aufklappbare Tipps. Tipp 3 enthält
  immer den kompletten Befehl zum Kopieren.
- **Irgendwas ist kaputt** → `bash scripts/ctf-reset.sh` räumt alles auf, dann von vorn.
- **Was läuft gerade überhaupt?** → `bash scripts/ctf-status.sh`
- **Spickzettel** → [docs/CHEATSHEET.md](docs/CHEATSHEET.md)

<details>
<summary>🚨 <strong>Komplettlösung</strong> (wirklich erst klicken, wenn du fertig bist!)</summary>

[docs/SOLUTION.de.md](docs/SOLUTION.de.md) — verrät alle Flags und alle Befehle.

</details>

---

<a name="english"></a>
## 🇬🇧 English

### What is this?

You work through five levels. Each level teaches you **one** new Docker command and
unlocks a **flag** — a short code like `FLAG{...}`. The big final flag waits at the end.

You do **not** need to write any code. You type commands into a terminal and look at the
result in your browser.

### What you need

- A machine with **Docker** (the bootcamp VM has everything already).
- A **terminal**.
- A **browser**.
- About **90 minutes**.

Check that Docker is running:

```bash
docker --version
```

### What you will learn

| Level | You learn |
|---|---|
| 0 | Start, inspect, stop and clean up containers |
| 1 | Look inside an image (`docker inspect`) and publish ports (`-p`) |
| 2 | Read metadata (labels) and name containers (`--name`) |
| 3 | Let containers talk to each other (`docker network`) |
| 4 | Share files between host and container (`-v`) |
| 5 | Run commands inside a running container (`docker exec -it`) |

### Get started

Pull both images first:

```bash
docker pull matixmedia/docker-ctf-main:latest
```

```bash
docker pull matixmedia/docker-ctf-data-provider:latest
```

Then open the workbook — it walks you through everything:

### ➡️ **[docs/WORKBOOK.en.md](docs/WORKBOOK.en.md)** ⬅️

> 💡 The CTF web interface is available in German and English. Switch in the top right.

### Help

- **I'm stuck** → Every level has three collapsible hints. Hint 3 always contains the
  complete command to copy.
- **Something broke** → `bash scripts/ctf-reset.sh` cleans up, then start over.
- **What is even running?** → `bash scripts/ctf-status.sh`
- **Cheat sheet** → [docs/CHEATSHEET.md](docs/CHEATSHEET.md)

<details>
<summary>🚨 <strong>Full solution</strong> (only click when you are done!)</summary>

[docs/SOLUTION.en.md](docs/SOLUTION.en.md) — reveals every flag and every command.

</details>

---

## 🧑‍🏫 Für Trainer / For trainers

- [docs/TRAINER.md](docs/TRAINER.md) — Ablauf, Timing, häufige Fehlerbilder, Pre-Flight-Check
- [docs/SLIDES-ADDENDUM.md](docs/SLIDES-ADDENDUM.md) — fehlende Folien für die Präsentation
- `docker compose up --build` — CTF lokal bauen und starten, ohne Docker Hub
- `bash scripts/build-and-push.sh` — Images bauen, taggen und pushen
