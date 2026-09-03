# Handoff #2: Weave the hands-on warm-up into the presentation

> **Internal working document.** Companion to
> [HANDOFF-PRESENTATION.md](HANDOFF-PRESENTATION.md) — read that one first for the deck
> inventory, the slide map and the three missing concept slides. This document covers one
> change only: **turning the Docker section into a follow-along session.**
>
> Every command below was executed and verified. Where a demo behaves differently from
> what you would expect, it is called out in §6 — **read that section before writing any
> slide**, two of the obvious demos silently do the wrong thing.

---

## 1. The change

Originally the plan had a separate **"Level 0"** warm-up block: a 15-minute guided nginx
exercise sitting between the deck and the CTF.

**Max wants it distributed into the presentation instead.** While he explains a concept at
the front, the students immediately try that exact command themselves. So the Docker
section stops being a lecture followed by an exercise, and becomes
*explain → everyone types → check the room → next concept*.

Your job: design that interleaving and the slides that carry it.

## 2. Why this is the right shape

- **It is what worked last year.** Deck page 40 ("ToDos fürs nächste Jahr") is the team's
  own retro: *"Hier hat eine Einführung in docker gefehlt, bei unserer Gruppe hat's ganz
  gut funktioniert dass wir es vorne gezeigt haben und alle mitgemacht haben."* The
  concrete suggestion on that slide — *"Dummy http Server starten, anschließend eine
  lokale HTML Datei in den Server mounten"* — is **exactly** the thread in §4. Following
  it closes the retro.
- **The CTF's biggest drop-off is level 1**, because it is the first time a participant
  types a docker command with nobody driving. If they have already typed `run`, `logs`,
  `ps`, `stop`, `rm`, `inspect`, `--network` and `-v` once with the room, level 1 stops
  being a cliff.
- **The audience needs it.** Mixed Azubi group — Fachinformatik, but also
  Wirtschaftsinformatiker and data scientists, several with no CLI experience at all.

## 3. Where each hands-on moment goes

PDF page / deck footer "Seite", per the map in handoff #1.

| After deck slide | Concept just taught | Hands-on moment |
|---|---|---|
| 36 / 37 — Docker Images verwalten | `pull`, `images` | **M1** pull nginx, list images |
| 37 / 38 — Einen Container starten | `run -d -p --name` | **M2** start nginx, open it in the browser |
| 38 / 39 — Mit einem Container interagieren | `exec -it`, `logs` | **M3** read logs, go inside the container |
| 39 / 40 — Weitere wichtige Docker Befehle | `ps`, `stop`, `rm` | **M4** lifecycle **+ provoke the name conflict on purpose** |
| **NEW** — `docker inspect` | inspect, labels | **M5** inspect a real image |
| **NEW** — Docker networks | networks, DNS | **M6** make it fail, then make it work |
| **NEW** — Volumes / bind mounts | `-v`, `$(pwd)` | **M7** serve your own HTML file |
| — | | cleanup, then the CTF slide |

The three **NEW** slides are already drafted in
[SLIDES-ADDENDUM.md](SLIDES-ADDENDUM.md) — restyle them to the template, then attach M5–M7.

> **Ordering note (changed from handoff #1).** Put the practice for the *basic* commands
> (M1–M4) directly after the slides that teach them, i.e. straight after page 39, and only
> then the three new concept slides with M5–M7. Sequence: basics → practise basics → new
> concepts → practise those → CTF.

## 4. The thread — one nginx container, built up step by step

Deliberately **one continuous story**, not seven disconnected snippets: the same container
gains a capability at each step, which is precisely the shape the CTF then asks for alone.

All commands verified working.

### M1 — an image is not a container
```bash
docker pull nginx
docker images
```
Point to make: nothing is running yet. An image is a template.

### M2 — your first container
```bash
docker run -d -p 8080:80 --name warmup nginx
```
Then everyone opens **<http://localhost:8080>** → *"Welcome to nginx!"*
This is the moment the room wakes up. Do not rush it.

### M3 — look inside
```bash
docker logs warmup
docker exec -it warmup sh
```
Inside the container:
```bash
ls /usr/share/nginx/html
exit
```
Ties straight back to the Linux/Bash part of the deck: same `ls`, same `cat`.

### M4 — lifecycle, and the error they *will* meet later
Provoke the name conflict **on purpose** — this is the single most valuable 60 seconds of
the whole warm-up:
```bash
docker ps
docker run -d --name warmup nginx      # <-- fails on purpose
```
Verified output:
```
docker: Error response from daemon: Conflict. The container name "/warmup"
is already in use by container "49b6189267dd...". You have to remove (or
rename) that container to be able to reuse that name.
```
Then the fix, which becomes the reflex they need all through the CTF:
```bash
docker stop warmup
docker ps -a          # still there, just stopped
docker rm warmup
```

### M5 — look at the metadata
```bash
docker run -d --name web nginx
docker inspect --format '{{.Config.ExposedPorts}}' web
```
Verified output: `map[80/tcp:{}]`

Say explicitly: *"the CTF's first level is exactly this, with a port you don't know yet."*

### M6 — why containers can't find each other
Keep `web` from M5 running (it is on the default bridge).
```bash
docker run --rm alpine wget -T2 -qO- http://web
```
Verified output — a clear, visible failure:
```
wget: can't connect to remote host (127.0.53.53): Connection refused
```
Now the fix:
```bash
docker rm -f web
docker network create demo-net
docker run -d --name web --network demo-net nginx
docker run --rm --network demo-net alpine wget -T2 -qO- http://web
```
→ the nginx welcome HTML. Name resolution works only on a self-created network.

> ⚠️ **Use `wget`, not `ping`.** See §6.1 — `ping` makes this demo *look like it succeeds*.

### M7 — your own file inside the container
The retro's own suggestion. The strongest demo in the set.
```bash
mkdir site
echo '<h1>Hallo Otto!</h1>' > site/index.html

docker rm -f web
docker run -d -p 8080:80 --name web --network demo-net \
  -v $(pwd)/site:/usr/share/nginx/html nginx
```
Browser → **Hallo Otto!**

Then edit `site/index.html`, reload the browser — it changes **immediately, without
restarting the container**. Verified. That is the moment bind mounts click.

### Cleanup before the CTF
```bash
docker rm -f web
docker network rm demo-net
```
Then the CTF slide. The CTF uses different names (`ctf-main`, `data-provider-svc`,
`ctf-net`) and port 8989, so nothing collides — but a clean start avoids confusion.

## 5. Slide design rules

- **Make follow-along slides visually distinct** from teaching slides at a glance. The deck
  already has a teal/mint accent (page 19, "Alles Bereit?") — reuse it, or a consistent
  `⌨️ MITMACHEN` / `⌨️ HANDS-ON` label.
- **Max 3 commands per slide.** People are typing, not reading.
- Every hands-on slide shows **"Das solltet ihr sehen:"** with the expected output. Without
  it, a participant cannot tell success from failure and will not speak up.
- End each one with an explicit sync point — *"Klappt's bei allen?"* — so Max can catch the
  stragglers before moving on. This is what stops the group silently splitting into
  "following" and "lost".
- M4 and M6 are **failure demos**. Mark them clearly as *"das soll jetzt schiefgehen"*,
  otherwise half the room thinks they broke something.
- Commands must be **copy-pasteable and typo-tolerant**: no smart quotes, no line-wrapped
  commands without a `\`.

## 6. Verified gotchas — read before writing slides

### 6.1 `ping` destroys the network demo 🔴
The obvious demo is `docker run --rm alpine ping -c1 web` on the default bridge. Verified
result:
```
PING web (127.0.53.53): 56 data bytes
64 bytes from 127.0.53.53: seq=0 ttl=64 time=0.018 ms
```
It **looks like it succeeded**. The resolver hands back `127.0.53.53` (the ICANN
name-collision address) and ping happily replies to itself. A participant would conclude
that containers *can* find each other by name — the exact opposite of the lesson.

Use `wget -T2 -qO- http://web`, which fails visibly. (Observed on Colima/macOS; on the
Linux bootcamp VM the message may instead be *"bad address"* — either way `wget` fails
where `ping` misleads. **Re-run this on the actual bootcamp VM** and put the real message
on the slide.)

### 6.2 Bind mounts only work from a shared path 🟠
`-v $(pwd)/site:/usr/share/nginx/html` from a directory under `/tmp` silently mounts an
**empty** folder, and nginx answers `403 Forbidden`. Verified. The Docker VM only shares
certain host paths (typically the home directory).

On the native Linux bootcamp VM this is not an issue, but tell participants to
`cd ~` first anyway — it costs nothing and removes a confusing failure. Same applies to
CTF level 4 (`-v $(pwd)/secrets:/secrets`).

### 6.3 Pre-pull the images 🟠
`nginx` (~50 MB) and `alpine` (~8 MB) plus the CTF image. Twenty people pulling
simultaneously over conference wifi is several minutes of dead air. Put it in the invite,
or run it in the break before.

### 6.4 Port 8080 must be free
Anyone with something already on 8080 gets `port is already allocated`. Have `-p 8081:80`
ready as the fallback answer.

### 6.5 `docker exec -it` needs a real terminal
Fine in a normal terminal; fails inside some IDE consoles. If a participant sees
*"the input device is not a TTY"*, they are not in a real shell.

## 7. Timing

| | |
|---|---|
| M1–M4 | ~2 min each → **~8 min** |
| M5 | ~3 min |
| M6 | ~6 min (two runs plus the explanation) |
| M7 | ~6 min |
| Sync/stragglers | ~5 min |
| **Total added to the Docker section** | **~28 min** |

This **replaces** the separate 15-minute Level 0 block, so the net cost is roughly
**+13 minutes** — inside a 2–3 h workshop, with the CTF still needing ~90 minutes.
If time gets tight, M5 is the one to shorten (it is repeated as CTF level 1 anyway).
**Do not cut M6 or M7** — those are the two concepts the deck is missing entirely.

## 8. Knock-on effects in this repo

- `docs/WORKBOOK.de.md` / `.en.md` still contain **Level 0 as a written section**. Once the
  deck owns the warm-up, that section becomes the **fallback** — for latecomers, absentees
  and anyone who wants to redo it at home. Keep it, but it is no longer the primary path.
- `docs/TRAINER.md` refers to "Level 0 gemeinsam" in its schedule but **does not contain
  the steps**. Once you have settled the slide sequence, that section should either carry
  the full script or point at the slides. Flag which you chose.
- **Nothing in the CTF itself changes.** Do not touch container names, image names, flags
  or level content.

## 9. Definition of done

- [ ] Every hands-on moment M1–M7 has a slide, visually distinct from teaching slides
- [ ] Each shows the expected output and ends with a sync point
- [ ] M4 and M6 are labelled as deliberate failures
- [ ] `wget` (not `ping`) is used in M6, with the message re-verified on the bootcamp VM
- [ ] `cd ~` appears before the first `-v` command
- [ ] A pre-pull instruction exists before the Docker section
- [ ] The three concept slides from `SLIDES-ADDENDUM.md` are in, before the CTF slide
- [ ] Total Docker section still leaves ~90 minutes for the CTF
- [ ] **Dry run**: type every command from the slides on a clean bootcamp VM, in order,
      and confirm each printed result matches what the slide claims
