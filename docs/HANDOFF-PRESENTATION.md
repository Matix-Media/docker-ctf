# Handoff: Improving the "Linux & Virtualisierung 2025" presentation

> **Internal working document.** This is a briefing for whoever (human or agent) picks up
> the presentation work. It is not participant material — feel free to delete it before
> sharing the repo with Azubis.
>
> Written after a full pass over the deck and a rebuild of the CTF in this repo.
> Everything below was verified against the actual PDF and the actual running containers,
> except where explicitly marked **[verify]**.

---

> 📎 **Companion document:** [HANDOFF-PRESENTATION-HANDSON.md](HANDOFF-PRESENTATION-HANDSON.md)
> covers weaving the hands-on warm-up into the Docker section — read it after §7 here.
> It supersedes the "separate Level 0 block" idea.

## 1. Your task

Improve the slide deck used in the Otto Linux-Bootcamp so that it actually prepares the
participants for the Docker CTF that follows it. The CTF side of this has already been
rebuilt (see §6) — **your job is the deck, not the CTF.**

The single most important outcome: **a participant who has seen the deck must be able to
start the CTF without meeting an unfamiliar command.** Right now that is not the case, and
it is the main reason last year's session went badly.

---

## 2. Source files

| What | Where |
|---|---|
| The deck (PDF, 48 pages, 2.6 MB) | `/Users/max/Downloads/Linux & Virtualisierung 2025.pdf` |
| Ready-made replacement/new slide text | `docs/SLIDES-ADDENDUM.md` (in this repo) |
| The CTF the deck leads into | this repo, start at `README.md` and `docs/WORKBOOK.de.md` |
| Trainer guide (timing, failure modes) | `docs/TRAINER.md` |

**Reading the PDF:** it is too large to read in one call. Use the `Read` tool with the
`pages` parameter, max 20 pages per request, e.g. `pages: "33-48"`. Reading without
`pages` will fail.

---

## 3. ⚠️ Blocking constraint — read this before you start

**The PDF is a rendering, not an editable deck.** You cannot produce a revised
presentation from it alone. Before doing any production work, ask Max:

1. **Where is the source file?** (`.pptx` / Keynote / Google Slides). Everything below
   assumes you get it. The Otto template, fonts (that condensed red display face) and
   colours live there.
2. **What deliverable does he want?**
   - edited source file returned, or
   - a written spec of slide changes he applies himself (this is what
     `docs/SLIDES-ADDENDUM.md` already is), or
   - a rebuilt deck from scratch.
3. If no source file exists, agree on the format *before* writing anything.

Do not silently produce a `.pptx` in a different visual style and call it the deck — the
Otto branding is not optional for an internal Otto training.

---

## 4. Background you need

**Audience.** Otto Azubis, mixed and much less technical than the deck assumes:
Fachinformatik Anwendungsentwicklung, but also **Wirtschaftsinformatiker and data
scientists**. Max's words: *"some even never used the cli"*. Assume genuine beginners.

**Format.** 2–3 hour workshop, trainer present, participants work in the Ubuntu
"Linux-Bootcamp-VM" (deck pages 18–20), with git and GitHub SSH already set up.

**Language.** Deck is German, informal *du*. Keep it that way. A few headings are English
(*"Good 2 Know"*, *"With great power comes great responsibility!"*, *"Checkout"*) — the
mixing is a bit arbitrary but it is a stylistic call, not a defect. Flag it, don't
unilaterally "fix" it.

**What went wrong last year.** The CTF was far too hard and the instructions were too
vague. The root cause turned out to be in the deck, not the CTF (see §5.1).

**The deck already contains last year's retro.** PDF page 40 ("ToDos fürs nächste Jahr")
says:

> *"Hier hat eine Einführung in docker gefehlt, bei unserer Gruppe hat's ganz gut
> funktioniert dass wir es vorne gezeigt haben und alle mitgemacht haben"*
> Was wir gemacht haben: Ubuntu im Docker ausführen und die Bash Befehle vom Anfang
> anwenden · Dummy http Server (httpd) starten · Anschließend eine lokale HTML Datei in
> den Server mounten

That is a strong signal and it has already been acted on in the CTF (Level 0 is a guided
warm-up everyone does together). **This slide is an internal note and should not be
presented to Azubis** — turn it into speaker notes or remove it once acted on.

---

## 5. Confirmed problems, ranked

### 5.1 The three missing concepts — the actual root cause 🔴

The CTF requires three tools that appear **nowhere in the 48-page deck**:

| Tool | Needed for | In the deck? |
|---|---|---|
| `docker inspect` (and `--format`) | CTF levels 1 and 2 | ❌ |
| `docker network create` / `--network` | CTF level 3 | ❌ |
| `-v` bind mounts | CTF level 4 | ❌ |

Deck pages 36–39 teach `pull`, `images`, `run` (`-d -p -e --name`), `exec -it`, `logs -f`,
`ps`, `ps -a`, `stop`, `rm`. Good, but it stops exactly one step short of what the CTF
needs. Participants hit level 1 and meet `docker inspect` for the first time.

**Reordering the slides does not fix this** — these concepts are absent from the entire
deck, not merely out of order. Three new slides are already written for you in
`docs/SLIDES-ADDENDUM.md`: `docker inspect`, Docker networks, volumes/bind mounts. They
are drafted in the deck's existing voice and structure (**Zweck** / **Beispiel** bullets).
Review them, restyle to the template, insert **before** the CTF slide.

### 5.2 CTF slide needs updating (but the image name is CORRECT) 🟠

PDF page 41 ("Selber probieren!") says:

> Docker Image ziehen: `matixmedia/docker-ctf`

**That name is right — do not change it.** `matixmedia/docker-ctf` is the main image
and exists on Docker Hub. (An earlier draft of this handoff claimed it was wrong; that
was a mistake. The repo's old `INSTRUCTIONS.md` named `docker-ctf-main`, which does *not*
exist — the repo was wrong, not the slide.)

What *does* need updating: the CTF now has **five levels with five flags**, and players
should be told to read the container logs. Replacement text is in
`docs/SLIDES-ADDENDUM.md` §"Anpassung 1".

> ⚠️ **Do not put a link to the source repository on the slide.** Players are given only
> the image name, by design: the repo contains the full solution, and the second image is
> meant to be discovered in-game.

### 5.3 Ordering 🟠

The CTF slide (page 41) currently sits **before** the Dockerfile/build slides (42–44) and
the Compose slides (45–47). **Max has already decided to reorder** so the entire Docker
teaching block comes before the CTF. Your job is to make that reordering coherent — the
agenda (page 3) and section dividers need to follow.

### 5.4 Dangling reference 🟡

Page 17 ("Good 2 Know") says **"+ Befehle auf Folie 69"**. The deck has 48 pages. Either
restore the intended slide or drop the reference.

### 5.5 Numbering gap — possible hidden slide 🟡

The deck's own footer numbers and the PDF page numbers diverge from page 32 onward:

- PDF page 31 → footer "Seite 31"
- PDF page 32 → footer **"Seite 33"**
- PDF page 33 → footer "Seite 34"

So a slide numbered 32 exists in the source but is not in the PDF — most likely a
**hidden slide**. Check the source deck for hidden slides; there may be content that was
meant to be there (possibly the missing Docker intro?).

### 5.6 Duplicate photo **[verify]** 🟡

Pages 6 (Timothy) and 7 (Tjark) appear to use the **same photograph**, and both carry the
identical bullets ("Ausbildungsjahrgang 2023", "FT3 = bestes Team"). Looks like a
copy-paste that was never finished. Verify against the source and ask the people involved
before changing anything — these are real colleagues' slides, not yours to rewrite.

### 5.7 Density and pacing 🟡

- Page 13 ("SHELL! und einige Befehle") fires ~14 commands in one list with no practice
  step. For someone who has never used a CLI this is where they start drowning.
- Page 46 (Compose YAML) is a dense config dump. Consider building it up in stages.
- The deck ends on page 47 (Compose commands) then straight to the OTTO end card — no
  wrap-up, no "what you learned", no links for going further.

### 5.8 External links to re-check 🟡

Page 27 points at `github.com/krother/bash_tutorial` and `github.com/veltman/clmystery`.
Verify both still exist and work before the session.

---

## 6. Already done — do not redo

The CTF in this repo has been fully rebuilt and verified end to end. **Do not change the
commands, flags, image names or level structure** — the deck must stay consistent with it.

- 5 levels, each with its own flag, plus a guided Level 0 warm-up
- Tiered hints (nudge → concept → complete copy-pasteable command)
- The web UI distinguishes "name does not resolve" from "nothing listening"
- Images are Alpine-based, so `docker exec -it ctf-main sh` works
- DE/EN toggle in the app; workbook and solution exist in both languages

Relevant files: `docs/WORKBOOK.de.md`, `docs/SOLUTION.de.md`, `docs/TRAINER.md`,
`docs/CHEATSHEET.md`.

---

## 7. The consistency contract

By the time the CTF slide appears, the deck must have taught **all** of this. Use it as
your acceptance checklist:

| Concept | Needed for | Currently taught? |
|---|---|---|
| `docker pull`, `docker images` | Setup | ✅ page 36 |
| `docker run -d -p --name` | Level 0/1 | ✅ page 37 |
| `docker ps`, `ps -a`, `stop`, `rm` | Level 0 + every restart | ✅ page 39 |
| `docker logs` | Level 1 | ✅ page 38 |
| **`docker inspect` + `--format`** | **Levels 1–2** | ❌ **add** |
| **Image labels** | **Level 2** | ❌ **add** |
| **`docker network create` / `--network`** | **Level 3** | ❌ **add** |
| **`-v` bind mount + `$(pwd)`** | **Level 4** | ❌ **add** |
| `docker exec -it` | Level 5 | ✅ page 38 |

One more worth adding as its own slide: **"a container name exists only once — `docker
stop` + `docker rm` before restarting"**. It is not a Docker concept so much as the single
most common way beginners get stuck. Draft text is in `docs/SLIDES-ADDENDUM.md`
("Optional: Folie D").

---

## 8. Complete slide map

PDF page numbers, with the deck's own footer number where it differs. Title/section/dark
slides carry no footer.

| PDF | "Seite" | Title |
|---|---|---|
| 1 | – | OTTO title card |
| 2 | – | Linux-Bootcamp 2025 — Was virtuelle Pinguine mit deiner Ausbildung zu tun haben |
| 3 | – | Agenda (6 points) |
| 4 | – | § Wer sind wir? |
| 5 | 5 | Max Heilmann |
| 6 | 6 | Timothy |
| 7 | 7 | Tjark **[see 5.6]** |
| 8 | 8 | Jan Moritz (JaMo) |
| 9 | 9 | § Linux |
| 10 | 10 | Linux – Einführung |
| 11 | 11 | Linux – Distributionen |
| 12 | 12 | Linux – Dateisystem |
| 13 | 13 | Linux – SHELL! und einige Befehle… **[dense, see 5.7]** |
| 14 | 14 | Linux – SHELL != SHELL |
| 15 | 15 | Bash Scripting: Automatisierung leicht gemacht |
| 16 | 16 | Bash Scripting: Aufbau eines Bash-Skripts |
| 17 | 17 | Good 2 Know – Einfacher mit Linux starten **[dangling ref, 5.4]** |
| 18 | 18 | § Linux-Bootcamp-VM |
| 19 | 19 | Alles Bereit? |
| 20 | 20 | GitHub SSH Key hinterlegen |
| 21 | 21 | § Maintenance |
| 22 | – | With great power comes great responsibility! |
| 23 | 23 | § Bash Escape Game |
| 24 | 24 | Umgebungsvariablen in der Bash-Shell |
| 25 | 25 | Programme ausführen in der Bash-Shell |
| 26 | 26 | Prozesse in der Bash-Shell |
| 27 | 27 | Checkout (bash_tutorial, clmystery) **[verify links, 5.8]** |
| 28 | 28 | § Virtualisierung |
| 29 | 29 | Virtualisierung – ein weit gefasster Oberbegriff |
| 30 | 30 | Beispiele aus der Otto-Welt |
| 31 | 31 | Was ist ein Hypervisor? |
| — | *32* | **missing / hidden — see 5.5** |
| 32 | 33 | § Docker Einführung |
| 33 | 34 | Warum Docker? |
| 34 | 35 | Docker packt alles in eine Box! |
| 35 | 36 | Container sind keine VMs! |
| 36 | 37 | Docker Images verwalten (`pull`, `images`) |
| 37 | 38 | Einen Container starten (`run -d -p -e --name`) |
| 38 | 39 | Mit einem Container interagieren (`exec -it`, `logs -f`) |
| 39 | 40 | Weitere wichtige Docker Befehle (`ps`, `ps -a`, `stop`, `rm`) |
| 40 | 41 | ToDos fürs nächste Jahr **[internal note, 5.4/§4]** |
| 41 | 42 | **Selber probieren! (CTF)** **[wrong image name, 5.2]** |
| 42 | 43 | Das Rezept für ein Docker Image (Dockerfile) |
| 43 | 44 | Ein einfaches Beispiel (Dockerfile) |
| 44 | 45 | Das eigene Docker Image bauen (`build -t`) |
| 45 | 46 | Mehrere Container einfach verwalten (Compose) |
| 46 | 47 | Web-App mit Datenbank (compose YAML) **[dense, 5.7]** |
| 47 | 48 | Die wichtigsten Compose Befehle |
| 48 | – | OTTO end card |

**Insertion point for the new slides:** after page 39 ("Weitere wichtige Docker Befehle")
and before the CTF slide.

---

## 9. Decisions already made — don't relitigate

- The Docker teaching block moves **entirely before** the CTF. Max is doing the reordering.
- The CTF stays as rebuilt: 5 levels, 5 flags, German + English in one image.
- Flags get posted in Teams so the trainer can see live who is stuck where. The CTF slide
  should say this.
- Participants work in the bootcamp VM, so `$(pwd)` and bash are safe — **no PowerShell
  variants needed** anywhere.

---

## 10. Ask Max about

1. **The source deck file** — blocking, see §3.
2. Slide 32: hidden on purpose, or lost content?
3. Timothy/Tjark slide — is the duplicate photo intentional?
4. Are the other presenters (Timothy, Tjark, JaMo) co-owners of the deck? Their sections
   may not be Max's to rewrite.
5. Does the Bash Escape Game section (23–27) stay? It affects the time budget: the CTF
   part needs ~90 minutes.
6. Should the Dockerfile/Compose slides (42–47) stay in at all? They are *after* the CTF
   in teaching order and the CTF does not need them — but they may be there deliberately
   as an outlook.

---

## 11. Definition of done

- [ ] Every concept in the §7 checklist is taught before the CTF slide
- [ ] The CTF slide names the correct image(s) and the five-flag/Teams format
- [ ] Page 17's "Folie 69" reference resolved
- [ ] "ToDos fürs nächste Jahr" is no longer a presented slide
- [ ] Slide ordering coherent, agenda and section dividers updated to match
- [ ] Otto template, fonts and colours intact
- [ ] Deck is consistent with `docs/WORKBOOK.de.md` — same image names, same commands,
      same container names (`ctf-main`, `data-provider-svc`), same network name (`ctf-net`)
- [ ] A dry run: read the revised deck start to finish, then read `docs/WORKBOOK.de.md`
      Level 1, and confirm no command appears there that the deck did not introduce
