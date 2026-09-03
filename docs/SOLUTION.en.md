# 🚨 Full Solution — Docker CTF

**Language:** English · [Deutsch](SOLUTION.de.md)

> ⚠️ **SPOILER WARNING**
> This contains **every flag** and **every command**. If you still want to solve the CTF,
> continue in the [workbook](WORKBOOK.en.md) instead — it has tiered hints, and hint 3
> already gives you the complete command for **one** level at a time.
>
> This file is mainly for **trainers** and for reading up **after** the workshop.

---

## Overview

| Level | Concept | Answer | Flag |
|---|---|---|---|
| 1 | `docker inspect`, `-p` | port `8989` | `FLAG{L1_P0RT_G3FUNDEN}` |
| 2 | Labels | `data-provider-svc` | `FLAG{L2_L4B3L_G3L3S3N}` |
| 3 | `docker network` | network `ctf-net` | `FLAG{L3_N3TZW3RK_ST3HT}` |
| 4 | bind mount `-v` | `SUPER_GEHEIM_123` | `FLAG{L4_TR3S0R_G3KN4CKT}` |
| 5 | `docker exec -it` | — | `FLAG{D0CK3R_PR0F1_MIT_FLAG}` |

---

## Preparation

```bash
docker pull matixmedia/docker-ctf-main:latest
docker pull matixmedia/docker-ctf-data-provider:latest
```

---

## Level 1 — The Secret Port

### Commands

```bash
docker run -d --name ctf-main matixmedia/docker-ctf-main:latest
docker logs ctf-main
```

**Expected output (shortened):**

```
  ==================================================================
    Docker CTF — Level 1: Der geheime Port / The Secret Port
  ==================================================================
...
  EN  I am listening on a secret port, but you cannot reach me from
      the outside yet. Find the port:

          docker inspect --format '{{.Config.ExposedPorts}}' ctf-main
...
Server laeuft. / Server is running.
```

```bash
docker inspect --format '{{.Config.ExposedPorts}}' ctf-main
```

**Output:**

```
map[8989/tcp:{}]
```

```bash
docker stop ctf-main
docker rm ctf-main

docker run -d --name ctf-main -p 8989:8989 matixmedia/docker-ctf-main:latest
```

Then open <http://localhost:8989> → **`FLAG{L1_P0RT_G3FUNDEN}`**

### Why does this work?

`EXPOSE 8989` in the Dockerfile is **documentation only**. It opens nothing. The port is
only forwarded from the host into the container by `-p 8989:8989` on `docker run`. The
left side is the port on your machine, the right side the one in the container — they do
not have to match.

`docker stop` + `docker rm` is required because a running container cannot be extended
with ports (or networks, or mounts) after the fact. You always create a new one.

---

## Level 2 — The Invisible Friend

### Commands

```bash
docker inspect --format '{{.Config.Labels}}' ctf-main
```

**Output:**

```
map[ctf.data-provider.host:data-provider-svc org.opencontainers.image.source:https://github.com/Matix-Media/docker-ctf org.opencontainers.image.title:Docker CTF - Main]
```

The value you need is **`data-provider-svc`**. Enter it on the web page
→ **`FLAG{L2_L4B3L_G3L3S3N}`**

### Why does this work?

`LABEL "ctf.data-provider.host"="data-provider-svc"` in the Dockerfile attaches arbitrary
metadata to the image. Labels change no behaviour — they are a place to store information
that belongs to the image (maintainer, version, source, and here the name of a related
service).

---

## Level 3 — The Network

### Commands

```bash
docker stop ctf-main data-provider-svc
docker rm ctf-main data-provider-svc

docker network create ctf-net

docker run -d --name data-provider-svc --network ctf-net \
  matixmedia/docker-ctf-data-provider:latest

docker run -d --name ctf-main --network ctf-net -p 8989:8989 \
  matixmedia/docker-ctf-main:latest
```

Reload the page → green box → **`FLAG{L3_N3TZW3RK_ST3HT}`**

### Why does this work?

On the default network (`bridge`) Docker runs **no** DNS for container names. Containers
can only reach each other by IP address there, and those can change on every restart.

A **self-created** network comes with a built-in DNS server. On it, a container's `--name`
becomes its hostname. The main container calls the hard-coded address
`http://data-provider-svc:9090/ping` — which only works if **both** containers are on the
same self-created network.

### The two failure modes

The web page distinguishes them explicitly:

| Message | Meaning | Cause |
|---|---|---|
| "I do not know that name" | DNS fails | container not on the same network |
| "I find the name, but nobody answers" | DNS works, TCP does not | container not running / wrong image |

---

## Level 4 — The Vault

### Commands

```bash
docker stop ctf-main
docker rm ctf-main

docker run -d --name ctf-main --network ctf-net -p 8989:8989 \
  -v $(pwd)/secrets:/secrets \
  matixmedia/docker-ctf-main:latest

cat secrets/password.txt
```

**Output:**

```
SUPER_GEHEIM_123
```

Enter the password on the web page → **`FLAG{L4_TR3S0R_G3KN4CKT}`**

### Why does this work?

At startup the application checks whether `/secrets` exists and only then writes the file.
Without a mount the folder does not exist in the image — so nothing happens. That is why
the container has to be restarted **after** adding the mount.

`-v $(pwd)/secrets:/secrets` is a **bind mount**: Docker maps the host folder directly into
the container. Both sides see the same files. Docker requires an absolute path here, hence
`$(pwd)`. If the host folder does not exist yet, Docker creates it.

The password is then verified over HTTP by the data provider (`POST /verify`) — that is,
over exactly the network connection from level 3.

---

## Level 5 — The Flag

### Command

```bash
docker exec -it ctf-main /app/app --show-flag
```

Then press **ENTER**.

**Output:**

```
Du bist fast am Ziel! Drücke ENTER, um die Flagge anzuzeigen.
You are almost there! Press ENTER to reveal the flag.

  FLAG{D0CK3R_PR0F1_MIT_FLAG}

Congratulations, you completed every level! 🐳
```

### Why does this work?

`docker exec` starts an **additional process** in an already running container. It is the
same program as the main process, but invoked with the argument `--show-flag` — so the
`CMD` from the Dockerfile is not replaced here, a second invocation simply runs alongside
it.

Without `-it` the program refuses to print anything:

- `docker exec ctf-main ...` → stdin is `/dev/null`
- `docker exec -i ctf-main ...` → stdin is a pipe, not a terminal
- `docker exec -it ctf-main ...` → stdin is a real pseudo terminal ✅

The program checks this with `ioctl(TCGETS)` — the classic "isatty" test. In the first two
cases the participant gets an error message that names the **complete** correct command.

### Bonus

```bash
docker exec -it ctf-main sh
```

The images are based on Alpine, so there is a shell. That lets participants explore the
file system with `ls`, `cd` and `cat` — exactly as in the Linux part of the bootcamp.

---

## Clean up

```bash
bash scripts/ctf-reset.sh
```

Also clear the cookies for `localhost:8989` so the progress bar resets too.
