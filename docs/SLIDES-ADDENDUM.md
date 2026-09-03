# Folien-Ergänzung für „Linux & Virtualisierung"

Diese Datei enthält fertigen Text für **drei fehlende Folien** sowie **zwei Korrekturen**
an der bestehenden Präsentation.

## Warum das nötig ist

Das CTF verlangt drei Werkzeuge, die im gesamten Deck (48 Seiten) **nicht vorkommen**:

| Werkzeug | Wird gebraucht in | Im Deck erklärt? |
|---|---|---|
| `docker inspect` | Level 1 und 2 | ❌ nein |
| `docker network create` / `--network` | Level 3 | ❌ nein |
| `-v` (Bind Mount / Volume) | Level 4 | ❌ nein |

Folien 37–40 lehren `pull`, `images`, `run` (`-d -p -e --name`), `exec -it`, `logs -f`,
`ps`, `ps -a`, `stop`, `rm`. Das reicht für Level 0 — aber ab Level 1 laufen die
Teilnehmer gegen ein Werkzeug, das sie nie gesehen haben.

> ⚠️ Das Umsortieren der Folien behebt das **nicht**. Diese drei Folien müssen neu dazu,
> und zwar **vor** der CTF-Folie.

---

## Korrektur 1 — Folie 42 „Selber probieren!"

**Falsch (aktuell):**

> Docker Image ziehen: `matixmedia/docker-ctf`

Dieses Image existiert nicht. Wer das eintippt, bekommt einen `pull access denied`-Fehler
und kommt gar nicht erst los.

**Richtig:**

> Docker Images ziehen:
> ```
> docker pull matixmedia/docker-ctf-main:latest
> docker pull matixmedia/docker-ctf-data-provider:latest
> ```
> Anleitung: `github.com/Matix-Media/docker-ctf` → `docs/WORKBOOK.de.md`
>
> Es gibt **fünf Level mit je einer eigenen Flag**. Postet eure Flags in Teams —
> dann sehen wir, wo ihr steht, und ihr verratet den anderen nichts.

## Korrektur 2 — Folie 17 „Good 2 Know"

Der Verweis **„+ Befehle auf Folie 69"** zeigt ins Leere — das Deck hat 48 Seiten.
Entweder die gemeinte Folie ergänzen oder den Verweis entfernen.

---

# Neue Folie A — In den Container hineinschauen

> **Titel:** In den Container hineinschauen: `docker inspect`

**Jedes Image und jeder Container bringt Metadaten mit**

Ein Image weiß mehr über sich, als man von außen sieht: welchen Port es benutzt, welchen
Startbefehl es hat, welche Zusatzinfos (Labels) die Entwickler hinterlegt haben.

`docker inspect <container oder image>`
- **Zweck:** Zeigt **alles**, was Docker über den Container weiß — als JSON.
- **Achtung:** Das sind gerne mal 200 Zeilen. Nicht erschrecken!

**Gezielt einen einzelnen Wert herausholen**

`docker inspect --format '{{.Config.ExposedPorts}}' <name>`
- **Zweck:** Statt der ganzen JSON-Wüste nur den einen Wert, den du brauchst.
- **Beispiel:** `docker inspect --format '{{.Config.ExposedPorts}}' mein-webserver`
  → `map[80/tcp:{}]`

`docker inspect --format '{{.Config.Labels}}' <name>`
- **Zweck:** Zeigt die **Labels** — frei vergebbare Zusatzinfos, die im Dockerfile stehen.
- **Beispiel-Ausgabe:** `map[maintainer:otto ctf.data-provider.host:data-provider-svc]`

> 💡 **Merke:** `EXPOSE` im Dockerfile *dokumentiert* nur einen Port.
> Erreichbar wird er erst mit `-p` beim `docker run`.

---

# Neue Folie B — Container reden miteinander

> **Titel:** Container reden miteinander: Docker-Netzwerke

**Das Problem**

Zwei Container laufen auf demselben Rechner. Trotzdem findet Container A den Container B
**nicht** unter seinem Namen. Warum?

Weil Docker alle Container standardmäßig in ein gemeinsames Standard-Netzwerk
(`bridge`) steckt — und **dort gibt es keine Namensauflösung**. Die Container haben nur
IP-Adressen, die sich bei jedem Neustart ändern können.

**Die Lösung: ein eigenes Netzwerk**

`docker network create <netzwerk-name>`
- **Zweck:** Erstellt ein eigenes Netzwerk mit **eingebautem DNS**.
- **Beispiel:** `docker network create mein-netz`

`docker run --network <netzwerk-name> ...`
- **Zweck:** Startet einen Container **in** diesem Netzwerk.
- **Beispiel:** `docker run -d --name datenbank --network mein-netz postgres`

`docker network ls`
- **Zweck:** Zeigt alle vorhandenen Netzwerke an.

**Das Merkbild**

> Im **eigenen** Netzwerk wird der `--name` eines Containers zu seinem **Hostnamen**.
> `--name datenbank` → andere Container im selben Netzwerk erreichen ihn unter
> `http://datenbank:5432`.
>
> **Beide** Container müssen im selben Netzwerk sein. Einer allein reicht nicht!

---

# Neue Folie C — Daten teilen

> **Titel:** Daten teilen: Volumes und Bind Mounts

**Das Problem**

Ein Container ist eine geschlossene Box. Löschst du ihn mit `docker rm`, sind alle
Dateien darin **weg**. Und von außen kommst du an sie gar nicht erst ran.

**Die Lösung: einen Ordner durchreichen**

`docker run -v <pfad-auf-deinem-rechner>:<pfad-im-container> ...`
- **Zweck:** Verbindet einen Ordner von deinem Rechner mit einem Ordner im Container.
  Beide sehen **dieselben** Dateien — in beide Richtungen, sofort.
- **Beispiel:** `docker run -d -v $(pwd)/website:/usr/share/nginx/html nginx`
  → Die HTML-Dateien aus deinem Ordner `website` werden vom nginx im Container ausgeliefert.

**Was ist `$(pwd)`?**
- `pwd` = *print working directory* = „in welchem Ordner stehe ich gerade?"
- `$(pwd)` setzt genau diesen Pfad ein. Docker braucht nämlich einen **kompletten** Pfad,
  kein `./website`.

> 💡 **Merke:** Zwei Anwendungsfälle —
> **rein:** eine Konfigurationsdatei in den Container geben.
> **raus:** Daten (Logs, Datenbank, Ergebnisse) behalten, auch wenn der Container weg ist.

---

# Optional: Folie D — Die goldene Aufräum-Regel

> Diese eine Regel spart im Workshop erfahrungsgemäß die meiste Zeit.

**Ein Containername kann nur einmal vergeben werden.**

Startest du einen Container mit `--name ctf-main` ein zweites Mal, kommt:

```
docker: Error response from daemon: Conflict. The container name "/ctf-main"
is already in use by container "a1b2c3...". You have to remove (or rename)
that container to be able to reuse that name.
```

**Das ist kein Fehler von dir.** Der alte Container existiert nur noch.
Vor **jedem** Neustart mit gleichem Namen also:

```bash
docker stop ctf-main && docker rm ctf-main
```

> 💡 Tipp: `docker run --rm ...` räumt den Container beim Stoppen automatisch weg.
