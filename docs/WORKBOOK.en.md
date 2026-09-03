# 🐳 Docker CTF — Workbook

**Language:** English · [Deutsch](WORKBOOK.de.md)

Welcome! You will solve five levels. Each one teaches you **one** new Docker command and
unlocks a **flag**.

> **You do not have to write any code.** You type commands into a terminal and check in the
> browser what happened.

## How to use this workbook

- Every grey box is a command. **You may copy and paste it.**
- Stuck? Every level has **three hints** you can unfold.
  Hint 3 always contains the **complete** command. That is not cheating.
- Write your flags down below and post them in Teams.

## 🚨 The single most important rule

A container name exists only **once**. To start a container with the same name again, you
have to remove the old one first:

```bash
docker stop ctf-main && docker rm ctf-main
```

If you ever see this:

```
docker: Error response from daemon: Conflict. The container name "/ctf-main"
is already in use ...
```

… you did **nothing wrong**. You just forgot to clean up. Run the two commands above and
try again.

---

# Level 0 — Warm-up 🔥

*No flag. We do this one together.* ⏱️ ~15 minutes

Before we start, let's practise the commands on a harmless web server.

### 0.1 Is Docker running?

```bash
docker --version
```

Do you get a version number? Good. If you get an error, say so.

### 0.2 Download an image

```bash
docker pull nginx
```

An **image** is a template — like a cake recipe. Nothing is running yet.

### 0.3 Start a container

```bash
docker run -d -p 8080:80 --name warmup nginx
```

Now a real web server is running. The three options:

| Option | Meaning |
|---|---|
| `-d` | *detached* — runs in the background, your terminal stays free |
| `-p 8080:80` | port **8080 on your machine** is forwarded to port **80 in the container** |
| `--name warmup` | gives the container a memorable name |

Now open **<http://localhost:8080>** — it says "Welcome to nginx!".

### 0.4 Have a look

```bash
docker ps
```

```bash
docker logs warmup
```

`docker ps` shows running containers, `docker logs` shows what a container prints.

### 0.5 Clean up — the routine you will need constantly

```bash
docker stop warmup
```

```bash
docker ps -a
```

The container is **stopped but still there** (`ps -a` also shows stopped ones). Only this
deletes it for good:

```bash
docker rm warmup
```

✅ **If all of that worked, you are ready for the CTF.**

---

# Level 1 — The Secret Port 🔍

⏱️ ~15 minutes

The CTF's main container is a web server — but it will not tell you which port it listens
on. You have to find that out.

### 📘 New for you: `docker inspect`

Every image carries **metadata**: which port it uses, which start command it has, which
notes the developers left in it. `docker inspect` shows you all of it — as JSON, and
easily 200 lines of it. With `--format` you pull out exactly **one** value.

### Your task

Start the container, read the logs, find the port and publish it.

```bash
docker run -d --name ctf-main matixmedia/docker-ctf-main:latest
```

```bash
docker logs ctf-main
```

The container tells you in its log what to do next. Read it carefully.

<details>
<summary>💡 Hint 1 — a small nudge</summary>

The log contains a command using `docker inspect`. Run it. It shows you a line like
`map[XXXX/tcp:{}]` — the number in there is the port you are looking for.

</details>

<details>
<summary>💡💡 Hint 2 — the concept behind it</summary>

Two steps:

1. Find the port:
   `docker inspect --format '{{.Config.ExposedPorts}}' ctf-main`
2. **Restart** the container and publish the port. An `EXPOSE` in the image only
   *documents* the port — it becomes reachable only with `-p PORT:PORT`.

Do not forget the cleanup (see the rule at the top).

</details>

<details>
<summary>💡💡💡 Hint 3 — the complete command</summary>

```bash
docker inspect --format '{{.Config.ExposedPorts}}' ctf-main
```

Output: `map[8989/tcp:{}]` → the port is **8989**.

```bash
docker stop ctf-main
docker rm ctf-main

docker run -d --name ctf-main -p 8989:8989 matixmedia/docker-ctf-main:latest
```

Then open it in your browser: <http://localhost:8989>

</details>

### ✅ Done when …

… you can open **<http://localhost:8989>** in your browser and see the CTF page.
Your first flag is on it.

**My flag:** `FLAG{ ________________________ }`

---

# Level 2 — The Invisible Friend 🏷️

⏱️ ~10 minutes

The page shows a red box: the main container needs a second container, the **data
provider**, but cannot reach it. First you have to find out what it is supposed to be
called.

### 📘 New for you: labels

An image can carry **labels**: little notes the developers wrote into it — for example who
built it or what a related service is called. You see them with `docker inspect` again.

### Your task

Find the name in the label `ctf.data-provider.host` and enter it on the web page.

<details>
<summary>💡 Hint 1 — a small nudge</summary>

You already know the command from level 1. This time you are not looking under
`Config.ExposedPorts` but under `Config.Labels`.

</details>

<details>
<summary>💡💡 Hint 2 — the concept behind it</summary>

```
docker inspect --format '{{.Config.Labels}}' YOUR_CONTAINER_NAME
```

In the output, look for the entry `ctf.data-provider.host:` — the value you need comes
right after it.

</details>

<details>
<summary>💡💡💡 Hint 3 — the complete command</summary>

```bash
docker inspect --format '{{.Config.Labels}}' ctf-main
```

The output contains, among other things:

```
ctf.data-provider.host:data-provider-svc
```

So the name is **`data-provider-svc`**. Enter it in the field on the web page.

</details>

### ✅ Done when …

… the web page accepts your answer and shows you the level 2 flag.

**My flag:** `FLAG{ ________________________ }`

---

# Level 3 — The Network 🔌

⏱️ ~20 minutes · *This is the trickiest level. Take your time.*

You now know what the data provider must be called. But even if you start it with the
right name, the main container will **not** find it.

### 📘 New for you: Docker networks

By default Docker puts all containers on one shared network — but **there is no name
resolution there**. The containers only have IP addresses, no names.

Only on a **self-created** network does a container's `--name` become its **hostname**.
Then container A reaches container B simply at `http://name-of-b`.

> ⚠️ **Both** containers have to be on the network. Putting only one on it is not enough —
> that is the mistake almost everyone makes.

### Your task

Create a network and start **both** containers on it.

The data provider is in the image `matixmedia/docker-ctf-data-provider:latest`.

> 💬 **Read the red box on the web page!** It now tells you exactly *what* is wrong:
> - *"I do not know that name"* → you are not on the same network.
> - *"I find the name, but nobody answers"* → the container is not running.

<details>
<summary>💡 Hint 1 — a small nudge</summary>

You need three things: a network, the data provider **on it** (with the name from level 2)
and the main container **on it as well**.

The command to create one is `docker network create`.

</details>

<details>
<summary>💡💡 Hint 2 — the concept behind it</summary>

1. `docker network create NETNAME`
2. Start the data provider with `--name data-provider-svc --network NETNAME`
3. **Restart** the main container, also with `--network NETNAME`

Step 3 is the one most often forgotten. The main container is already running — but on the
wrong network. A running container cannot simply be moved, you have to recreate it.

</details>

<details>
<summary>💡💡💡 Hint 3 — the complete command</summary>

Copy the whole block. The first two lines clean up — you can ignore any errors they
produce.

```bash
docker stop ctf-main data-provider-svc
docker rm ctf-main data-provider-svc

docker network create ctf-net

docker run -d --name data-provider-svc --network ctf-net \
  matixmedia/docker-ctf-data-provider:latest

docker run -d --name ctf-main --network ctf-net -p 8989:8989 \
  matixmedia/docker-ctf-main:latest
```

Then reload <http://localhost:8989>.

</details>

### ✅ Done when …

… the red box on the web page turns **green**.

**My flag:** `FLAG{ ________________________ }`

---

# Level 4 — The Vault 🔐

⏱️ ~15 minutes

The main container knows a password. But it only writes it into a folder **inside** its
container: `/secrets/password.txt`. You cannot reach that from outside.

### 📘 New for you: bind mounts (`-v`)

A container is a closed box — delete it and every file inside is gone.
With `-v` you connect a folder **on your machine** to a folder **in the container**. Both
then see the same files, in both directions, instantly.

```
-v YOUR_FOLDER:/FOLDER_IN_CONTAINER
```

Docker wants a **full** path here. `./secrets` is not enough.
`$(pwd)` automatically inserts the folder you are currently in.

### Your task

Mount a folder to `/secrets`, restart the main container, read the password and enter it on
the web page.

<details>
<summary>💡 Hint 1 — a small nudge</summary>

The container only writes the file if the folder `/secrets` already exists **at startup**.
So you have to hand it in **and then restart**.

Careful: the main container has to stay on the `ctf-net` network, otherwise you fall back
to level 3!

</details>

<details>
<summary>💡💡 Hint 2 — the concept behind it</summary>

Take your `docker run` command from level 3 and append one option:

```
-v $(pwd)/secrets:/secrets
```

Afterwards you have a folder `secrets` containing `password.txt`. You can read it with
`cat`.

</details>

<details>
<summary>💡💡💡 Hint 3 — the complete command</summary>

```bash
docker stop ctf-main
docker rm ctf-main

docker run -d --name ctf-main --network ctf-net -p 8989:8989 \
  -v $(pwd)/secrets:/secrets \
  matixmedia/docker-ctf-main:latest

cat secrets/password.txt
```

The last command shows you the password. Enter it on the web page.

</details>

### ✅ Done when …

… you have the file `secrets/password.txt` and the web page accepts the password.

**My flag:** `FLAG{ ________________________ }`

---

# Level 5 — The Flag 🏁

⏱️ ~10 minutes

The container only shows the final flag **inside itself** — and only if you bring a real
terminal.

### 📘 New for you: `docker exec -it`

With `docker exec` you run a command **inside an already running** container. You do not
have to build anything new.

The `-it` is two things:
- `-i` = *interactive* — your keystrokes reach the container
- `-t` = *TTY* — the container gets a real terminal

Without `-it` the program refuses to print anything. Try it, then you will see the
difference.

### Your task

Run the command `/app/app --show-flag` inside the main container.

<details>
<summary>💡 Hint 1 — a small nudge</summary>

The command starts with `docker exec`, then the container name, then the command that
should run inside it.

</details>

<details>
<summary>💡💡 Hint 2 — the concept behind it</summary>

```
docker exec -it CONTAINERNAME /app/app --show-flag
```

If you leave out the `-it` you get an error message instead of the flag — but the message
tells you exactly what is missing.

</details>

<details>
<summary>💡💡💡 Hint 3 — the complete command</summary>

```bash
docker exec -it ctf-main /app/app --show-flag
```

Then press **ENTER** once. 🎉

</details>

### 🎁 Bonus

Feel free to look around inside the container:

```bash
docker exec -it ctf-main sh
```

With `ls`, `cd` and `cat` you move through the file system — exactly the commands from the
Linux part. `exit` gets you back out.

### ✅ Done when …

… the final `FLAG{…}` appears in your terminal.

**My final flag:** `FLAG{ ________________________ }`

---

# 🎉 Done!

You inspected images, published ports, named containers, built a network, mounted a volume
and executed a command inside a running container. That is exactly what working with
Docker looks like day to day.

### Clean up

```bash
bash scripts/ctf-reset.sh
```

### When nothing works any more

| Problem | Solution |
|---|---|
| "container name is already in use" | `docker stop NAME && docker rm NAME` |
| I do not know what is running | `bash scripts/ctf-status.sh` |
| Everything is broken, start over | `bash scripts/ctf-reset.sh` |
| Page does not load | Is the container running? `docker ps`. Port published? `docker port ctf-main` |
| Progress bar looks wrong | Clear the cookies for `localhost:8989` or use a private window |

📄 Cheat sheet with all commands: [CHEATSHEET.md](CHEATSHEET.md)
