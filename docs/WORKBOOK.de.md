# 🐳 Docker CTF — Workbook

**Sprache:** Deutsch · [English](WORKBOOK.en.md)

Willkommen! Du löst hier fünf Level. In jedem lernst du **einen** neuen Docker-Befehl und
schaltest eine **Flag** frei.

> **Du musst nichts programmieren.** Du tippst Befehle in ein Terminal und schaust im
> Browser nach, was passiert ist.

## So benutzt du dieses Workbook

- Jeder graue Kasten ist ein Befehl. **Du darfst ihn kopieren und einfügen.**
- Kommst du nicht weiter? Jedes Level hat **drei Tipps** zum Aufklappen.
  Tipp 3 enthält immer den **kompletten** Befehl. Das ist kein Schummeln.
- Trag deine Flags unten ein und poste sie in Teams.

## 🚨 Die wichtigste Regel überhaupt

Einen Containernamen gibt es nur **einmal**. Willst du einen Container mit demselben Namen
neu starten, musst du den alten vorher wegräumen:

```bash
docker stop ctf-main && docker rm ctf-main
```

Wenn du jemals das hier siehst:

```
docker: Error response from daemon: Conflict. The container name "/ctf-main"
is already in use ...
```

… dann hast du **nichts falsch gemacht**. Du hast nur das Aufräumen vergessen. Führe die
beiden Befehle oben aus und probier es nochmal.

---

# Level 0 — Aufwärmen 🔥

*Ohne Flag. Das machen wir gemeinsam.* ⏱️ ~15 Minuten

Bevor es losgeht, üben wir die Befehle einmal an einem harmlosen Webserver.

### 0.1 Läuft Docker?

```bash
docker --version
```

Kommt eine Versionsnummer? Gut. Kommt ein Fehler, sag Bescheid.

### 0.2 Ein Image herunterladen

```bash
docker pull nginx
```

Ein **Image** ist eine Vorlage — wie ein Kuchenrezept. Noch läuft nichts.

### 0.3 Einen Container starten

```bash
docker run -d -p 8080:80 --name warmup nginx
```

Jetzt läuft ein echter Webserver. Die drei Optionen:

| Option | Bedeutung |
|---|---|
| `-d` | *detached* — läuft im Hintergrund, dein Terminal bleibt frei |
| `-p 8080:80` | Port **8080 bei dir** wird auf Port **80 im Container** geleitet |
| `--name warmup` | gibt dem Container einen merkbaren Namen |

Öffne jetzt **<http://localhost:8080>** — da steht „Welcome to nginx!".

### 0.4 Nachschauen

```bash
docker ps
```

```bash
docker logs warmup
```

`docker ps` zeigt laufende Container, `docker logs` zeigt, was ein Container ausgibt.

### 0.5 Aufräumen — die Routine, die du gleich ständig brauchst

```bash
docker stop warmup
```

```bash
docker ps -a
```

Der Container ist **gestoppt, aber noch da** (`ps -a` zeigt auch gestoppte). Erst das hier
löscht ihn endgültig:

```bash
docker rm warmup
```

✅ **Wenn das alles geklappt hat, bist du bereit für das CTF.**

---

# Level 1 — Der geheime Port 🔍

⏱️ ~15 Minuten

Der Hauptcontainer des CTF ist ein Webserver — aber er verrät dir nicht, auf welchem Port
er lauscht. Den musst du herausfinden.

### 📘 Das ist neu für dich: `docker inspect`

Jedes Image bringt **Metadaten** mit: welchen Port es benutzt, welchen Startbefehl es hat,
welche Notizen die Entwickler hinterlegt haben. `docker inspect` zeigt dir das alles — als
JSON, und zwar gerne mal 200 Zeilen. Mit `--format` holst du dir gezielt **einen** Wert
heraus.

### Deine Aufgabe

Starte den Container, lies die Logs, finde den Port und gib ihn frei.

```bash
docker run -d --name ctf-main matixmedia/docker-ctf-main:latest
```

```bash
docker logs ctf-main
```

Der Container sagt dir im Log, was du als Nächstes tun sollst. Lies es in Ruhe.

<details>
<summary>💡 Tipp 1 — ein kleiner Stups</summary>

Im Log steht ein Befehl mit `docker inspect`. Führ ihn aus. Er zeigt dir eine Zeile wie
`map[XXXX/tcp:{}]` — die Zahl darin ist der gesuchte Port.

</details>

<details>
<summary>💡💡 Tipp 2 — das Konzept dahinter</summary>

Zwei Schritte:

1. Port herausfinden:
   `docker inspect --format '{{.Config.ExposedPorts}}' ctf-main`
2. Container **neu starten** und den Port freigeben. Ein `EXPOSE` im Image
   *dokumentiert* den Port nur — erreichbar wird er erst mit `-p PORT:PORT`.

Vergiss das Aufräumen nicht (siehe die Regel ganz oben).

</details>

<details>
<summary>💡💡💡 Tipp 3 — der komplette Befehl</summary>

```bash
docker inspect --format '{{.Config.ExposedPorts}}' ctf-main
```

Ausgabe: `map[8989/tcp:{}]` → der Port ist **8989**.

```bash
docker stop ctf-main
docker rm ctf-main

docker run -d --name ctf-main -p 8989:8989 matixmedia/docker-ctf-main:latest
```

Dann im Browser öffnen: <http://localhost:8989>

</details>

### ✅ Geschafft, wenn …

… du **<http://localhost:8989>** im Browser öffnen kannst und die CTF-Seite siehst.
Dort steht deine erste Flag.

**Meine Flag:** `FLAG{ ________________________ }`

---

# Level 2 — Der unsichtbare Freund 🏷️

⏱️ ~10 Minuten

Die Seite zeigt einen roten Kasten: Der Hauptcontainer braucht einen zweiten Container,
den **Data-Provider**, erreicht ihn aber nicht. Erst mal musst du herausfinden, wie der
überhaupt heißen soll.

### 📘 Das ist neu für dich: Labels

Ein Image kann **Labels** tragen: kleine Notizzettel, die die Entwickler hineingeschrieben
haben — zum Beispiel wer es gebaut hat oder wie ein zugehöriger Dienst heißt. Du siehst sie
wieder mit `docker inspect`.

### Deine Aufgabe

Finde im Label `ctf.data-provider.host` den Namen und trage ihn auf der Webseite ein.

<details>
<summary>💡 Tipp 1 — ein kleiner Stups</summary>

Du kennst den Befehl schon aus Level 1. Diesmal schaust du nicht unter
`Config.ExposedPorts`, sondern unter `Config.Labels`.

</details>

<details>
<summary>💡💡 Tipp 2 — das Konzept dahinter</summary>

```
docker inspect --format '{{.Config.Labels}}' NAME_DEINES_CONTAINERS
```

In der Ausgabe suchst du den Eintrag `ctf.data-provider.host:` — direkt dahinter steht der
Wert, den du brauchst.

</details>

<details>
<summary>💡💡💡 Tipp 3 — der komplette Befehl</summary>

```bash
docker inspect --format '{{.Config.Labels}}' ctf-main
```

In der Ausgabe steht unter anderem:

```
ctf.data-provider.host:data-provider-svc
```

Der Name ist also **`data-provider-svc`**. Trag ihn auf der Webseite in das Feld ein.

</details>

### ✅ Geschafft, wenn …

… die Webseite deine Antwort annimmt und dir die Flag für Level 2 zeigt.

**Meine Flag:** `FLAG{ ________________________ }`

---

# Level 3 — Das Netzwerk 🔌

⏱️ ~20 Minuten · *Das ist das kniffligste Level. Nimm dir Zeit.*

Du weißt jetzt, wie der Data-Provider heißen muss. Aber selbst wenn du ihn mit dem
richtigen Namen startest, findet der Hauptcontainer ihn **nicht**.

### 📘 Das ist neu für dich: Docker-Netzwerke

Docker steckt alle Container standardmäßig in ein gemeinsames Netzwerk — aber **dort gibt
es keine Namensauflösung**. Die Container haben nur IP-Adressen, keine Namen.

Erst in einem **selbst erstellten** Netzwerk wird der `--name` eines Containers zu seinem
**Hostnamen**. Dann erreicht Container A den Container B einfach unter `http://name-von-b`.

> ⚠️ **Beide** Container müssen im Netzwerk sein. Nur einen hineinzustecken reicht nicht —
> das ist der Fehler, den fast alle machen.

### Deine Aufgabe

Erstelle ein Netzwerk und starte **beide** Container darin.

Der Data-Provider steckt im Image `matixmedia/docker-ctf-data-provider:latest`.

> 💬 **Lies den roten Kasten auf der Webseite!** Er sagt dir jetzt genau, *woran* es
> gerade hakt:
> - *„Ich kenne den Namen nicht"* → ihr seid nicht im selben Netzwerk.
> - *„Ich finde den Namen, aber niemand macht auf"* → der Container läuft nicht.

<details>
<summary>💡 Tipp 1 — ein kleiner Stups</summary>

Du brauchst drei Dinge: ein Netzwerk, den Data-Provider **darin** (mit dem Namen aus
Level 2) und den Hauptcontainer **ebenfalls darin**.

Der Befehl zum Anlegen heißt `docker network create`.

</details>

<details>
<summary>💡💡 Tipp 2 — das Konzept dahinter</summary>

1. `docker network create NETZNAME`
2. Data-Provider starten mit `--name data-provider-svc --network NETZNAME`
3. Hauptcontainer **neu starten**, ebenfalls mit `--network NETZNAME`

Schritt 3 wird am häufigsten vergessen. Der Hauptcontainer läuft ja schon — aber im
falschen Netzwerk. Ein laufender Container lässt sich nicht einfach umhängen, du musst ihn
neu erzeugen.

</details>

<details>
<summary>💡💡💡 Tipp 3 — der komplette Befehl</summary>

Kopiere den ganzen Block. Die ersten beiden Zeilen räumen auf — Fehlermeldungen dabei
kannst du ignorieren.

```bash
docker stop ctf-main data-provider-svc
docker rm ctf-main data-provider-svc

docker network create ctf-net

docker run -d --name data-provider-svc --network ctf-net \
  matixmedia/docker-ctf-data-provider:latest

docker run -d --name ctf-main --network ctf-net -p 8989:8989 \
  matixmedia/docker-ctf-main:latest
```

Danach die Seite <http://localhost:8989> neu laden.

</details>

### ✅ Geschafft, wenn …

… der rote Kasten auf der Webseite **grün** wird.

**Meine Flag:** `FLAG{ ________________________ }`

---

# Level 4 — Der Tresor 🔐

⏱️ ~15 Minuten

Der Hauptcontainer kennt ein Passwort. Er legt es aber nur in einem Ordner **innerhalb**
seines Containers ab: `/secrets/password.txt`. Da kommst du von außen nicht ran.

### 📘 Das ist neu für dich: Bind Mounts (`-v`)

Ein Container ist eine geschlossene Box — löschst du ihn, sind alle Dateien darin weg.
Mit `-v` verbindest du einen Ordner **von deinem Rechner** mit einem Ordner **im
Container**. Beide sehen dann dieselben Dateien, in beide Richtungen, sofort.

```
-v ORDNER_BEI_DIR:/ORDNER_IM_CONTAINER
```

Docker will dabei einen **vollständigen** Pfad. `./secrets` reicht nicht.
`$(pwd)` setzt automatisch den Ordner ein, in dem du gerade stehst.

### Deine Aufgabe

Mounte einen Ordner nach `/secrets`, starte den Hauptcontainer neu, lies das Passwort und
trage es auf der Webseite ein.

<details>
<summary>💡 Tipp 1 — ein kleiner Stups</summary>

Der Container schreibt die Datei nur, wenn der Ordner `/secrets` beim **Start** schon
existiert. Du musst ihn also hineinreichen **und danach neu starten**.

Achtung: Der Hauptcontainer muss weiterhin im Netzwerk `ctf-net` bleiben, sonst fällst du
auf Level 3 zurück!

</details>

<details>
<summary>💡💡 Tipp 2 — das Konzept dahinter</summary>

Nimm deinen `docker run`-Befehl aus Level 3 und häng eine Option an:

```
-v $(pwd)/secrets:/secrets
```

Danach liegt bei dir ein Ordner `secrets` mit der Datei `password.txt` darin. Auslesen
kannst du sie mit `cat`.

</details>

<details>
<summary>💡💡💡 Tipp 3 — der komplette Befehl</summary>

```bash
docker stop ctf-main
docker rm ctf-main

docker run -d --name ctf-main --network ctf-net -p 8989:8989 \
  -v $(pwd)/secrets:/secrets \
  matixmedia/docker-ctf-main:latest

cat secrets/password.txt
```

Der letzte Befehl zeigt dir das Passwort. Trag es auf der Webseite ein.

</details>

### ✅ Geschafft, wenn …

… bei dir die Datei `secrets/password.txt` liegt und die Webseite das Passwort annimmt.

**Meine Flag:** `FLAG{ ________________________ }`

---

# Level 5 — Die Flagge 🏁

⏱️ ~10 Minuten

Die finale Flag zeigt der Container nur **direkt bei sich drin** an — und nur, wenn du ein
echtes Terminal mitbringst.

### 📘 Das ist neu für dich: `docker exec -it`

Mit `docker exec` führst du einen Befehl **in einem bereits laufenden** Container aus.
Du musst dafür nichts neu bauen.

Das `-it` besteht aus zwei Teilen:
- `-i` = *interactive* — deine Tastatureingaben kommen im Container an
- `-t` = *TTY* — der Container bekommt ein echtes Terminal

Ohne `-it` verweigert das Programm die Ausgabe. Probier es ruhig aus, dann siehst du den
Unterschied.

### Deine Aufgabe

Führe im Hauptcontainer den Befehl `/app/app --show-flag` aus.

<details>
<summary>💡 Tipp 1 — ein kleiner Stups</summary>

Der Befehl fängt mit `docker exec` an, dann kommt der Containername, dann der Befehl, der
darin laufen soll.

</details>

<details>
<summary>💡💡 Tipp 2 — das Konzept dahinter</summary>

```
docker exec -it CONTAINERNAME /app/app --show-flag
```

Wenn du das `-it` weglässt, bekommst du eine Fehlermeldung statt der Flag — die Meldung
sagt dir aber genau, was fehlt.

</details>

<details>
<summary>💡💡💡 Tipp 3 — der komplette Befehl</summary>

```bash
docker exec -it ctf-main /app/app --show-flag
```

Dann einmal **ENTER** drücken. 🎉

</details>

### 🎁 Bonus

Schau dich ruhig im Container um:

```bash
docker exec -it ctf-main sh
```

Mit `ls`, `cd` und `cat` bewegst du dich durch das Dateisystem — genau die Befehle aus dem
Linux-Teil. Mit `exit` kommst du wieder heraus.

### ✅ Geschafft, wenn …

… im Terminal die finale `FLAG{…}` steht.

**Meine finale Flag:** `FLAG{ ________________________ }`

---

# 🎉 Fertig!

Du hast Images inspiziert, Ports freigegeben, Container benannt, ein Netzwerk gebaut, ein
Volume gemountet und einen Befehl in einem laufenden Container ausgeführt. Genau das macht
man im Alltag mit Docker.

### Aufräumen

```bash
bash scripts/ctf-reset.sh
```

### Wenn mal gar nichts mehr geht

| Problem | Lösung |
|---|---|
| „container name is already in use" | `docker stop NAME && docker rm NAME` |
| Ich weiß nicht, was gerade läuft | `bash scripts/ctf-status.sh` |
| Alles kaputt, neu anfangen | `bash scripts/ctf-reset.sh` |
| Seite lädt nicht | Läuft der Container? `docker ps`. Port veröffentlicht? `docker port ctf-main` |
| Fortschrittsbalken stimmt nicht | Cookies für `localhost:8989` löschen oder privates Fenster benutzen |

📄 Spickzettel mit allen Befehlen: [CHEATSHEET.md](CHEATSHEET.md)
