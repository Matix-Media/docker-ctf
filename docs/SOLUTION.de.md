# 🚨 Komplettlösung — Docker CTF

**Sprache:** Deutsch · [English](SOLUTION.en.md)

> ⚠️ **SPOILER-WARNUNG**
> Hier stehen **alle Flags** und **alle Befehle**. Wenn du das CTF noch lösen willst,
> mach lieber im [Workbook](WORKBOOK.de.md) weiter — dort gibt es gestufte Tipps, und
> Tipp 3 verrät dir immer schon den kompletten Befehl für **ein** Level.
>
> Diese Datei ist vor allem für **Trainer** und zum Nachlesen **nach** dem Workshop.

---

## Übersicht

| Level | Konzept | Antwort | Flag |
|---|---|---|---|
| 1 | `docker inspect`, `-p` | Port `8989` | `FLAG{L1_P0RT_G3FUNDEN}` |
| 2 | Labels | `data-provider-svc` | `FLAG{L2_L4B3L_G3L3S3N}` |
| 3 | `docker network` | Netzwerk `ctf-net` | `FLAG{L3_N3TZW3RK_ST3HT}` |
| 4 | Bind Mount `-v` | `SUPER_GEHEIM_123` | `FLAG{L4_TR3S0R_G3KN4CKT}` |
| 5 | `docker exec -it` | — | `FLAG{D0CK3R_PR0F1_MIT_FLAG}` |

---

## Vorbereitung

```bash
docker pull matixmedia/docker-ctf-main:latest
docker pull matixmedia/docker-ctf-data-provider:latest
```

---

## Level 1 — Der geheime Port

### Befehle

```bash
docker run -d --name ctf-main matixmedia/docker-ctf-main:latest
docker logs ctf-main
```

**Erwartete Ausgabe (gekürzt):**

```
  ==================================================================
    Docker CTF — Level 1: Der geheime Port / The Secret Port
  ==================================================================

  DE  Ich lausche auf einem geheimen Port, aber von aussen kommst du
      noch nicht an mich heran. Finde den Port heraus:

          docker inspect --format '{{.Config.ExposedPorts}}' ctf-main
...
Server laeuft. / Server is running.
```

```bash
docker inspect --format '{{.Config.ExposedPorts}}' ctf-main
```

**Ausgabe:**

```
map[8989/tcp:{}]
```

```bash
docker stop ctf-main
docker rm ctf-main

docker run -d --name ctf-main -p 8989:8989 matixmedia/docker-ctf-main:latest
```

Dann <http://localhost:8989> öffnen → **`FLAG{L1_P0RT_G3FUNDEN}`**

### Warum funktioniert das?

`EXPOSE 8989` im Dockerfile ist **nur Dokumentation**. Es öffnet nichts. Der Port wird
erst durch `-p 8989:8989` beim `docker run` vom Host in den Container geleitet. Links
steht der Port auf deinem Rechner, rechts der im Container — sie müssen nicht gleich sein.

Das `docker stop` + `docker rm` ist nötig, weil ein laufender Container sich nicht
nachträglich um Ports (oder Netzwerke, oder Mounts) erweitern lässt. Man erzeugt immer
einen neuen.

---

## Level 2 — Der unsichtbare Freund

### Befehle

```bash
docker inspect --format '{{.Config.Labels}}' ctf-main
```

**Ausgabe:**

```
map[ctf.data-provider.host:data-provider-svc org.opencontainers.image.source:https://github.com/Matix-Media/docker-ctf org.opencontainers.image.title:Docker CTF - Main]
```

Der gesuchte Wert ist **`data-provider-svc`**. Auf der Webseite eintragen
→ **`FLAG{L2_L4B3L_G3L3S3N}`**

### Warum funktioniert das?

`LABEL "ctf.data-provider.host"="data-provider-svc"` im Dockerfile hängt beliebige
Metadaten an das Image. Labels ändern am Verhalten nichts — sie sind ein Ort, an dem man
Informationen ablegt, die zum Image gehören (Maintainer, Version, Quelle, hier eben der
Name eines zugehörigen Dienstes).

---

## Level 3 — Das Netzwerk

### Befehle

```bash
docker stop ctf-main data-provider-svc
docker rm ctf-main data-provider-svc

docker network create ctf-net

docker run -d --name data-provider-svc --network ctf-net \
  matixmedia/docker-ctf-data-provider:latest

docker run -d --name ctf-main --network ctf-net -p 8989:8989 \
  matixmedia/docker-ctf-main:latest
```

Seite neu laden → grüner Kasten → **`FLAG{L3_N3TZW3RK_ST3HT}`**

### Warum funktioniert das?

Im Standard-Netzwerk (`bridge`) betreibt Docker **kein** DNS für Containernamen. Container
erreichen sich dort nur über IP-Adressen, die sich bei jedem Neustart ändern können.

Ein **selbst erstelltes** Netzwerk bringt einen eingebauten DNS-Server mit. Dort wird der
`--name` eines Containers zu seinem Hostnamen. Der Hauptcontainer ruft fest verdrahtet
`http://data-provider-svc:9090/ping` auf — das funktioniert nur, wenn **beide** Container
im selben selbst erstellten Netzwerk sind.

### Die zwei Fehlerbilder

Die Webseite unterscheidet sie explizit:

| Meldung | Bedeutung | Ursache |
|---|---|---|
| „Ich kenne den Namen nicht" | DNS scheitert | Container nicht im selben Netzwerk |
| „Ich finde den Namen, aber niemand macht auf" | DNS klappt, TCP nicht | Container läuft nicht / falsches Image |

---

## Level 4 — Der Tresor

### Befehle

```bash
docker stop ctf-main
docker rm ctf-main

docker run -d --name ctf-main --network ctf-net -p 8989:8989 \
  -v $(pwd)/secrets:/secrets \
  matixmedia/docker-ctf-main:latest

cat secrets/password.txt
```

**Ausgabe:**

```
SUPER_GEHEIM_123
```

Passwort auf der Webseite eintragen → **`FLAG{L4_TR3S0R_G3KN4CKT}`**

### Warum funktioniert das?

Die Anwendung prüft beim Start, ob `/secrets` existiert, und schreibt die Datei nur dann.
Ohne Mount existiert der Ordner im Image nicht — also passiert nichts. Deshalb muss der
Container **nach** dem Hinzufügen des Mounts neu gestartet werden.

`-v $(pwd)/secrets:/secrets` ist ein **Bind Mount**: Docker hängt den Host-Ordner direkt in
den Container. Beide Seiten sehen dieselben Dateien. Docker verlangt hier einen absoluten
Pfad, deshalb `$(pwd)`. Existiert der Host-Ordner noch nicht, legt Docker ihn an.

Das Passwort wird dann per HTTP vom Data-Provider geprüft (`POST /verify`) — also über
genau die Netzwerkverbindung aus Level 3.

---

## Level 5 — Die Flagge

### Befehl

```bash
docker exec -it ctf-main /app/app --show-flag
```

Danach **ENTER** drücken.

**Ausgabe:**

```
Du bist fast am Ziel! Drücke ENTER, um die Flagge anzuzeigen.
You are almost there! Press ENTER to reveal the flag.

  FLAG{D0CK3R_PR0F1_MIT_FLAG}

Glückwunsch, du hast alle Level geschafft! 🐳
```

### Warum funktioniert das?

`docker exec` startet einen **zusätzlichen Prozess** in einem bereits laufenden Container.
Das Programm ist dasselbe wie der Hauptprozess, wird aber mit dem Argument `--show-flag`
aufgerufen — der `CMD` aus dem Dockerfile wird hier also nicht ersetzt, sondern es läuft
schlicht ein zweiter Aufruf daneben.

Ohne `-it` verweigert das Programm die Ausgabe:

- `docker exec ctf-main ...` → stdin ist `/dev/null`
- `docker exec -i ctf-main ...` → stdin ist eine Pipe, kein Terminal
- `docker exec -it ctf-main ...` → stdin ist ein echtes Pseudo-Terminal ✅

Das Programm prüft das per `ioctl(TCGETS)` — dem klassischen „isatty"-Test. In den ersten
beiden Fällen bekommt der Teilnehmer eine Fehlermeldung, die den **kompletten** richtigen
Befehl nennt.

### Bonus

```bash
docker exec -it ctf-main sh
```

Die Images basieren auf Alpine, es gibt also eine Shell. Damit lässt sich das Dateisystem
mit `ls`, `cd` und `cat` erkunden — genau wie im Linux-Teil des Bootcamps.

---

## Aufräumen

```bash
bash scripts/ctf-reset.sh
```

Zusätzlich die Cookies für `localhost:8989` löschen, damit auch der Fortschrittsbalken
zurückgesetzt ist.
