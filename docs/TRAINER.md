# 🧑‍🏫 Trainer-Guide — Docker CTF

Alles, was du zum Durchführen brauchst. Für den Otto Linux-Bootcamp-Workshop.

---

## Was letztes Jahr schiefging (und was sich geändert hat)

| Problem | Ursache | Behoben durch |
|---|---|---|
| „viel zu schwer" | Das CTF verlangte `docker inspect`, `docker network` und `-v` — **keins davon kam in der Präsentation vor** | 3 neue Folien: [SLIDES-ADDENDUM.md](SLIDES-ADDENDUM.md) |
| Widersprüchliche Image-Namen im Repo | Das alte `INSTRUCTIONS.md` nannte `docker-ctf-main`; dieses Image existiert **nicht**. Richtig ist `matixmedia/docker-ctf` (so wie auf Folie 42) | Einheitlich `matixmedia/docker-ctf` |
| Nach 60 min nichts vorzuweisen | Eine einzige Flag ganz am Ende | 5 Level mit je einer eigenen Flag |
| `docker exec -it ... bash` schlug fehl | Images waren `FROM scratch` → keine Shell, obwohl Folie 39 genau das lehrt | Images basieren jetzt auf Alpine |
| „Ich weiß nicht, was ich falsch mache" | Jeder Verbindungsfehler zeigte denselben Text | Seite unterscheidet DNS-Fehler von „Container läuft nicht" |
| Hänger an Name-Konflikten | Nirgends stand, dass man vorher `docker rm` braucht | Aufräum-Regel in Level 0, im Workbook und auf jeder Seite |

---

## 🔑 Grundprinzip: Spieler bekommen nur den Image-Namen

Die Teilnehmer erfahren **ausschließlich**:

```
matixmedia/docker-ctf
```

Kein Repo-Link, keine Anleitung, kein zweiter Image-Name. Alles Weitere steht in den
Container-Logs (Level 1) und danach in der Weboberfläche (Level 2–5). Das ist Absicht:

- Im Repo liegt die **Komplettlösung**.
- Das zweite Image sollen sie in **Level 3 selbst entdecken**.

Wenn du das Workbook austeilen willst: als **PDF oder Ausdruck**, nicht als Repo-Link.
Nötig ist es nicht — das CTF ist ohne lösbar.

---

## ✅ Pre-Flight-Check (am Vortag!)

```bash
# 1. Images sind aktuell und öffentlich erreichbar
docker pull matixmedia/docker-ctf:latest
docker pull matixmedia/docker-ctf-data-provider:latest

# 2. Kompletter Durchlauf auf einer frischen VM
bash scripts/ctf-reset.sh
# ... Level 1-5 aus docs/SOLUTION.de.md durchspielen

# 3. Aufräumen
bash scripts/ctf-reset.sh
```

**Vor Ort:**

- [ ] Folien umsortiert: **kompletter Docker-Teil vor** dem CTF
- [ ] Die drei neuen Folien aus [SLIDES-ADDENDUM.md](SLIDES-ADDENDUM.md) sind eingebaut
- [ ] Image-Name auf der CTF-Folie korrigiert
- [ ] **Teilnehmer lassen das Image vorab ziehen** (`docker pull matixmedia/docker-ctf`) —
      20 Leute, die gleichzeitig pullen, sind mehrere Minuten Wartezeit. Am besten schon
      in der Pause davor. Das zweite Image ziehen sie später selbst in Level 3.
- [ ] **Kein Repo-Link** an der Wand oder im Chat — dort steht die Lösung
- [ ] Image-Name steht sichtbar im Chat: `matixmedia/docker-ctf`
- [ ] Du hast [SOLUTION.de.md](SOLUTION.de.md) offen

---

## ⏱️ Ablauf (~90 Minuten CTF-Teil)

| Zeit | Was | Wie |
|---|---|---|
| 0:00–0:15 | **Level 0 gemeinsam** | Du machst es vorne vor, alle tippen mit. Nicht überspringen! |
| 0:15–0:30 | Level 1 | Erstes Level allein. Hier ist die Abbruchgefahr am größten. |
| 0:30–0:40 | Level 2 | Geht meist schnell, weil `inspect` schon sitzt. |
| 0:40–1:00 | Level 3 | **Das schwerste Level.** Plan mehr Zeit ein. |
| 1:00–1:15 | Level 4 | |
| 1:15–1:25 | Level 5 | |
| 1:25–1:30 | Abschluss | Flags einsammeln, kurz durchsprechen, was passiert ist |

> 💡 **Level 0 ist die wichtigste Maßnahme.** Auf Folie 41 steht euer eigenes Feedback vom
> letzten Jahr: *„hat ganz gut funktioniert, dass wir es vorne gezeigt haben und alle
> mitgemacht haben."* Genau das ist Level 0.

### Flags einsammeln

Teilnehmer posten ihre Flags in Teams. So siehst du live, wer wo steht — und wenn nach
40 Minuten die halbe Gruppe noch keine `FLAG{L3_...}` gepostet hat, weißt du, dass du
Level 3 nochmal gemeinsam durchgehen solltest.

| Flag | Level |
|---|---|
| `FLAG{L1_P0RT_G3FUNDEN}` | 1 |
| `FLAG{L2_L4B3L_G3L3S3N}` | 2 |
| `FLAG{L3_N3TZW3RK_ST3HT}` | 3 |
| `FLAG{L4_TR3S0R_G3KN4CKT}` | 4 |
| `FLAG{D0CK3R_PR0F1_MIT_FLAG}` | 5 (final) |

---

## 🆘 Häufige Fehlerbilder

### „The container name is already in use"

```
docker: Error response from daemon: Conflict. The container name "/ctf-main"
is already in use by container "a1b2c3..."
```

**Das ist der mit Abstand häufigste Hänger.** Der alte Container existiert noch.

```bash
docker rm -f ctf-main
```

### Das Terminal hängt / reagiert nicht mehr

`-d` vergessen. Der Container läuft im Vordergrund. `Strg+C`, dann mit `-d` neu starten.

### Die Seite lädt nicht

Der Reihe nach prüfen:

```bash
docker ps                # läuft ctf-main überhaupt?
docker port ctf-main     # ist 8989 veröffentlicht?
docker logs ctf-main     # gibt es Fehler?
```

Oder einfach: `bash scripts/ctf-status.sh`

### „Ich habe den Data-Provider doch gestartet!"

Zu 90 % ist der **Hauptcontainer** nicht im Netzwerk. Beide müssen rein. Der rote Kasten
auf der Seite sagt es jetzt explizit:

- *„Ich kenne den Namen nicht"* → Netzwerkproblem (Level 3)
- *„Ich finde den Namen, aber niemand macht auf"* → Container läuft nicht (Level 2)

### `password.txt` erscheint nicht

Der Container schreibt die Datei nur **beim Start**. Wer den Mount hinzufügt, ohne den
Container neu zu erzeugen, wartet vergeblich.

### Der Fortschrittsbalken ist zurückgesprungen

Der Fortschritt liegt im Browser-Cookie (nicht im Container — der wird ja ständig neu
erzeugt). Anderer Browser, privates Fenster oder gelöschte Cookies = Fortschritt weg.
Die Flags bleiben aber gültig, und die Seite erkennt den echten Systemzustand von selbst
wieder.

### Jemand ist völlig durcheinander

```bash
bash scripts/ctf-reset.sh
```

Räumt Container, Netzwerk und `./secrets` weg. Danach bei Level 1 neu anfangen.

---

## 🎚️ Hinweise dosieren

Jedes Level hat drei Stufen — im Workbook **und** auf der Webseite:

| Stufe | Inhalt | Wann rausgeben |
|---|---|---|
| Tipp 1 | Stups in die richtige Richtung | nach ~5 min Grübeln |
| Tipp 2 | Das Konzept + Befehlsform ohne Werte | nach ~10 min |
| Tipp 3 | **Kompletter Befehl zum Kopieren** | wenn Frust aufkommt |

**Tipp 3 ist ausdrücklich erlaubt.** Lieber jemand kopiert den Befehl, sieht das Ergebnis
und versteht *danach* das Konzept, als dass er 20 Minuten feststeckt und abschaltet. Sag
das zu Beginn laut — sonst trauen sich die Unsicheren nicht.

Für die Schnellen: Sie sollen die Tipps zuklappen lassen und Level 5 mit dem Bonus
(`docker exec -it ctf-main sh`) erkunden, oder sich das Dockerfile im Repo ansehen.

---

## 🔧 Images neu bauen

```bash
# lokal testen, ohne Docker Hub
docker compose up --build

# bauen, taggen und pushen
bash scripts/build-and-push.sh

# nur bauen
PUSH=0 bash scripts/build-and-push.sh

# eigener Docker-Hub-Account / neues Jahr
REGISTRY=deinname VERSION=2026 bash scripts/build-and-push.sh
```

> `docker-compose.yml` ist **nur für dich**. Es konfiguriert Netzwerk, Namen und Volume
> bereits fertig — damit wären genau die Aufgaben gelöst, die die Teilnehmer selbst machen
> sollen. Nicht weitergeben.

Die Images werden zusätzlich versioniert getaggt (`:2025`), damit ein altes `latest` auf
einem Laptop den Workshop nicht sabotiert.

> ⚠️ **Multi-Arch ist Pflicht.** Auf einem Apple-Silicon-Mac erzeugt ein normales
> `docker build` **nur** `linux/arm64`. Die Teilnehmer-VMs sind aber `linux/amd64` — die
> Images liessen sich dort nicht starten. `scripts/build-and-push.sh` benutzt deshalb
> `docker buildx` und veröffentlicht beide Plattformen. Nach dem Push einmal gegenprüfen:
>
> ```bash
> docker manifest inspect matixmedia/docker-ctf:latest | grep architecture
> ```
>
> Es müssen `amd64` **und** `arm64` auftauchen.

> ⚠️ **Lokal testen:** Solange die neuen Images noch nicht gepusht sind, zieht ein
> `docker pull` (und auch `docker run` ohne lokales Image) die **alte** Version von Docker
> Hub und überschreibt damit deinen lokalen Build. Zum Testen also immer erst
> `PUSH=0 bash scripts/build-and-push.sh` und danach **kein** `docker pull`.

---

## 📁 Was liegt wo

| Datei | Zweck |
|---|---|
| `README.md` | Einstieg für Teilnehmer, zweisprachig |
| `docs/WORKBOOK.de.md` / `.en.md` | Das Arbeitsheft — hier arbeiten die Teilnehmer |
| `docs/SOLUTION.de.md` / `.en.md` | Komplettlösung mit Erklärungen (für dich) |
| `docs/CHEATSHEET.md` | Einseitiger Spickzettel, zweisprachig |
| `docs/SLIDES-ADDENDUM.md` | Die drei fehlenden Folien + Korrekturen |
| `scripts/ctf-reset.sh` | Alles aufräumen |
| `scripts/ctf-status.sh` | Selbstdiagnose: was läuft, was fehlt, was kommt als Nächstes |
| `scripts/build-and-push.sh` | Images bauen und veröffentlichen |
| `main-app/` | Hauptcontainer (Level 1–5) |
| `data-provider-app/` | Zweiter Container (Level 3–4) |
