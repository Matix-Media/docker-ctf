package main

import "html/template"

// Lang ist die Sprache der Oberfläche. / Lang is the UI language.
type Lang string

const (
	LangDE Lang = "de"
	LangEN Lang = "en"
)

// ParseLang macht aus einer beliebigen Eingabe eine gültige Sprache (Default: Deutsch).
func ParseLang(s string) Lang {
	if s == string(LangEN) {
		return LangEN
	}
	return LangDE
}

// Other gibt die jeweils andere Sprache zurück (für den Umschalter).
func (l Lang) Other() Lang {
	if l == LangDE {
		return LangEN
	}
	return LangDE
}

// T liefert einen übersetzten Text. Fehlt der Schlüssel, fällt er auf Deutsch zurück.
func T(l Lang, key string) string {
	if v, ok := translations[l][key]; ok {
		return v
	}
	if v, ok := translations[LangDE][key]; ok {
		return v
	}
	return key
}

// TH liefert einen übersetzten Text, der HTML enthalten darf.
func TH(l Lang, key string) template.HTML {
	return template.HTML(T(l, key))
}

var translations = map[Lang]map[string]string{

	// ---------------------------------------------------------------- DEUTSCH
	LangDE: {
		"page.title":     "Docker CTF",
		"brand":          "🐳 Docker Capture the Flag",
		"lang.switch":    "English",
		"progress.label": "Dein Fortschritt",
		"level.word":     "Level",

		"hints.label":        "Hinweise — klick auf, wenn du nicht weiterkommst",
		"conn.ok":            "Ich kann den Data-Provider erreichen. Die Verbindung steht!",
		"conn.bad":           "Ich erreiche den Data-Provider nicht.",
		"conn.anon.title":    "Data-Provider",
		"conn.anon.bad":      "Ich erreiche meinen Data-Provider nicht — und wie er heißt, verrate ich dir hier nicht. Genau das ist deine Aufgabe. 😉",
		"footer.resetcmd.l2": "docker rm -f ctf-main",
		"hint.1":             "💡 Tipp 1 — ein kleiner Stups",
		"hint.2":             "💡💡 Tipp 2 — das Konzept dahinter",
		"hint.3":             "💡💡💡 Tipp 3 — der komplette Befehl",

		"goal.label":    "Dein Ziel",
		"concept.label": "Neu für dich",
		"check.label":   "Woran du merkst, dass es geklappt hat",

		"flag.unlocked": "Flag freigeschaltet!",
		"flag.note":     "Schreib die Flag auf und poste sie in Teams — dann sehen wir, wie weit du bist.",

		"cleanup.title": "Die goldene Aufräum-Regel",
		"cleanup.body": `Einen Containernamen gibt es nur <strong>einmal</strong>. Bevor du einen Container
			mit demselben Namen neu startest, musst du den alten wegräumen. Sonst kommt
			<em>„The container name is already in use"</em> — das ist kein Fehler von dir!`,

		"footer.progress": "Dein Fortschritt wird als Cookie in diesem Browser gespeichert — nicht im Container. Ein privates Fenster fängt also wieder bei vorne an.",
		"footer.reset":    "Komplett verheddert? Damit räumst du alles ab und fängst bei Level 1 neu an:",
		"footer.resetcmd": "docker rm -f ctf-main data-provider-svc\ndocker network rm ctf-net",

		// ---- Level 1 (bereits geschafft, wenn diese Seite lädt)
		"l1.title": "Der geheime Port",
		"l1.done": `Du hast den versteckten Port gefunden und nach außen freigegeben. Genau deshalb
			siehst du diese Seite überhaupt. <strong>Level 1 geschafft!</strong>`,
		"l1.learned": `Du kannst jetzt mit <code>docker inspect</code> in einen Container hineinschauen und
			mit <code>-p</code> einen Port nach außen freigeben.`,

		// ---- Level 2
		"l2.title": "Der unsichtbare Freund",
		"l2.goal": `Ich brauche einen zweiten Container, den <strong>Data-Provider</strong>. Aber ich
			erreiche ihn nicht. Finde heraus, unter welchem <strong>Namen</strong> ich ihn suche,
			und starte ihn genau so.`,
		"l2.concept": `Ein Image kann <strong>Labels</strong> tragen: kleine Notizzettel, die die Entwickler
			hineingeschrieben haben. Genau so ein Zettel klebt an mir und verrät, wie mein Freund
			heißen muss. Du siehst Labels mit <code>docker inspect</code>.`,
		"l2.hint1": `Ich habe Metadaten, die du dir anschauen kannst — so wie du in Level 1 den Port
			gefunden hast. Diesmal suchst du nicht nach einem Port, sondern nach einem
			<strong>Label</strong>. Es heißt <code>ctf.data-provider.host</code>.`,
		"l2.hint2": `Mit <code>docker inspect</code> siehst du alles über einen Container — das sind aber
			schnell 200 Zeilen JSON. Mit <code>--format</code> holst du dir gezielt einen Wert heraus:
			<br><code>docker inspect --format '{{.Config.Labels}}' NAME_DEINES_CONTAINERS</code>`,
		"l2.hint3": `Führe das hier aus (dein Container heißt vermutlich <code>ctf-main</code>):
			<pre><code>docker inspect --format '{{.Config.Labels}}' ctf-main</code></pre>
			In der Ausgabe steht <code>ctf.data-provider.host:</code> und dahinter der gesuchte Name.
			Trage ihn unten ein.`,
		"l2.check":       `Wenn du den richtigen Namen unten einträgst, bekommst du die Flag für Level 2.`,
		"l2.question":    "Wie muss mein Freund, der Data-Provider, heißen?",
		"l2.placeholder": "Name aus dem Label eintragen …",
		"l2.submit":      "Antwort prüfen",
		"l2.wrong": `Das ist noch nicht der richtige Name. Schau dir die Labels an —
			der Schlüssel heißt <code>ctf.data-provider.host</code>.`,
		"l2.flagtext": "Du hast das Label gefunden und weißt jetzt, wie der Data-Provider heißen muss.",

		// ---- Level 3
		"l3.title": "Das Netzwerk",
		"l3.goal": `Der Data-Provider muss laufen <strong>und</strong> ich muss ihn erreichen können.
			Sorge dafür, dass wir beide im selben Docker-Netzwerk sind.
			<br><br>Sein Image findest du in der Docker Registry unter
			<code>matixmedia/docker-ctf-data-provider:latest</code>.`,
		"l3.concept": `Docker steckt alle Container standardmäßig in ein gemeinsames Netzwerk — aber
			<strong>dort gibt es keine Namensauflösung</strong>. Container finden sich also nicht über
			ihren Namen. Erst in einem <strong>selbst erstellten</strong> Netzwerk wird der
			<code>--name</code> eines Containers zu seinem Hostnamen.`,
		"l3.hint1": `Der Data-Provider ist im Image
			<code>matixmedia/docker-ctf-data-provider:latest</code>. Er muss laufen, den richtigen
			Namen tragen — und wir beide müssen uns im selben Netzwerk befinden.`,
		"l3.hint2": `Du brauchst drei Bausteine:
			<br>1. Ein Netzwerk anlegen: <code>docker network create NETZNAME</code>
			<br>2. Den Data-Provider darin starten: <code>--name data-provider-svc --network NETZNAME</code>
			<br>3. <strong>Mich neu starten</strong>, ebenfalls mit <code>--network NETZNAME</code>.
			<br>Punkt 3 wird gern vergessen — einer allein im Netzwerk reicht nicht!`,
		"l3.hint3": `Kopiere diesen Block komplett. Die ersten beiden Zeilen räumen auf,
			falls schon etwas läuft — Fehlermeldungen dabei kannst du ignorieren.
			<pre><code>docker stop ctf-main data-provider-svc
docker rm ctf-main data-provider-svc

docker network create ctf-net

docker run -d --name data-provider-svc --network ctf-net \
  matixmedia/docker-ctf-data-provider:latest

docker run -d --name ctf-main --network ctf-net -p 8989:8989 \
  matixmedia/docker-ctf:latest</code></pre>
			Danach diese Seite neu laden.`,
		"l3.check": `Sobald es klappt, wird der rote Kasten hier oben grün und du bekommst die Flag.`,
		"l3.diag.dns": `<strong>Ich kenne den Namen nicht.</strong> Der Name
			<code>data-provider-svc</code> lässt sich nicht auflösen. Das heißt: wir sind
			<strong>nicht im selben Netzwerk</strong>. Genau darum geht es in diesem Level.`,
		"l3.diag.refused": `<strong>Ich finde den Namen, aber niemand macht auf.</strong> Der Name wird
			aufgelöst — der Container läuft also nicht, ist abgestürzt, oder hört auf einem anderen
			Port. Prüfe mit <code>docker ps</code>, ob <code>data-provider-svc</code> wirklich läuft.`,
		"l3.flagtext": "Beide Container sind im selben Netzwerk und finden sich über ihren Namen.",

		// ---- Level 4
		"l4.title": "Der Tresor",
		"l4.goal": `Ich kenne ein Passwort, aber ich lege es nur in einem Ordner <em>innerhalb</em>
			meines Containers ab: <code>/secrets/password.txt</code>. Hol es dir auf deinen Rechner
			und trage es unten ein.`,
		"l4.concept": `Ein Container ist eine geschlossene Box — löschst du ihn, sind alle Dateien darin
			weg. Mit einem <strong>Bind Mount</strong> (<code>-v</code>) verbindest du einen Ordner von
			deinem Rechner mit einem Ordner im Container. Beide sehen dann <strong>dieselben</strong>
			Dateien.`,
		"l4.hint1": `Ich schreibe die Datei nur dann, wenn der Ordner <code>/secrets</code> auch
			wirklich existiert. Du musst ihn mir also von außen hineinreichen — und danach musst du
			mich <strong>neu starten</strong>, denn geschrieben wird beim Start.`,
		"l4.hint2": `Die Option heißt <code>-v ORDNER_BEI_DIR:/secrets</code>.
			Docker will dabei einen <strong>vollständigen</strong> Pfad — <code>./secrets</code> reicht
			nicht. <code>$(pwd)</code> setzt automatisch den Pfad ein, in dem du gerade stehst.`,
		"l4.hint3": `Kopiere diesen Block komplett:
			<pre><code>docker stop ctf-main
docker rm ctf-main

docker run -d --name ctf-main --network ctf-net -p 8989:8989 \
  -v $(pwd)/secrets:/secrets \
  matixmedia/docker-ctf:latest

cat secrets/password.txt</code></pre>
			Der letzte Befehl zeigt dir das Passwort. Trage es unten ein.`,
		"l4.check": `Nach dem Neustart liegt bei dir ein Ordner <code>secrets</code> mit der Datei
			<code>password.txt</code> darin.`,
		"l4.mounted.no":  `❌ Ich sehe noch keinen <code>/secrets</code>-Ordner. Der Mount fehlt noch.`,
		"l4.mounted.yes": `✅ Der Ordner <code>/secrets</code> ist da — ich habe das Passwort hineingeschrieben. Lies es aus und trage es unten ein.`,
		"l4.pwlabel":     "Passwort",
		"l4.placeholder": "Passwort aus password.txt …",
		"l4.submit":      "Passwort prüfen",
		"l4.wrong": `Das Passwort stimmt nicht. Lies die Datei nochmal genau aus:
			<code>cat secrets/password.txt</code> — ohne Leerzeichen davor oder dahinter.`,
		"l4.flagtext": "Du hast einen Ordner in den Container gemountet und das Passwort geborgen.",

		// ---- Level 5
		"l5.title": "Die Flagge",
		"l5.goal": `Die finale Flag zeige ich nur direkt in meinem Container an — und nur, wenn du
			ein <strong>echtes Terminal</strong> mitbringst.`,
		"l5.concept": `Mit <code>docker exec</code> führst du einen Befehl <em>in</em> einem bereits
			laufenden Container aus. Das <code>-it</code> steht für <em>interactive</em> +
			<em>TTY</em>: Damit bekommt der Befehl eine echte Tastatur — ohne das geht es hier nicht.`,
		"l5.hint1": `Du musst mir keinen neuen Container bauen. Ich laufe schon. Du musst nur einen
			Befehl <strong>in mir</strong> ausführen: <code>/app/app --show-flag</code>`,
		"l5.hint2": `Der Befehl lautet <code>docker exec ...</code> — aber ohne
			<code>-it</code> verweigere ich die Ausgabe. Probier ruhig aus, was ohne passiert.`,
		"l5.hint3": `Der komplette Befehl:
			<pre><code>docker exec -it ctf-main /app/app --show-flag</code></pre>
			Danach einmal <strong>ENTER</strong> drücken — und die Flag gehört dir. 🎉`,
		"l5.check": `Im Terminal erscheint die finale <code>FLAG{…}</code>.`,
		"l5.bonus": `<strong>Bonus:</strong> Schau dich ruhig mal in mir um:
			<code>docker exec -it ctf-main sh</code> — mit <code>ls</code> und <code>cat</code> kannst
			du dich durch mein Dateisystem bewegen. Mit <code>exit</code> kommst du wieder raus.`,

		// ---- Fertig
		"done.title": "Alle Level geschafft!",
		"done.body": `Du hast Images inspiziert, Ports freigegeben, Container benannt, ein Netzwerk
			gebaut, ein Volume gemountet und einen Befehl in einem laufenden Container ausgeführt.
			Das ist genau das, was man im Alltag mit Docker macht. 🐳`,
	},

	// ---------------------------------------------------------------- ENGLISH
	LangEN: {
		"page.title":     "Docker CTF",
		"brand":          "🐳 Docker Capture the Flag",
		"lang.switch":    "Deutsch",
		"progress.label": "Your progress",
		"level.word":     "Level",

		"hints.label":        "Hints — open one if you get stuck",
		"conn.ok":            "I can reach the data provider. The connection is up!",
		"conn.bad":           "I cannot reach the data provider.",
		"conn.anon.title":    "Data provider",
		"conn.anon.bad":      "I cannot reach my data provider — and I am not telling you its name here. That is exactly your task. 😉",
		"footer.resetcmd.l2": "docker rm -f ctf-main",
		"hint.1":             "💡 Hint 1 — a small nudge",
		"hint.2":             "💡💡 Hint 2 — the concept behind it",
		"hint.3":             "💡💡💡 Hint 3 — the complete command",

		"goal.label":    "Your goal",
		"concept.label": "New for you",
		"check.label":   "How you know it worked",

		"flag.unlocked": "Flag unlocked!",
		"flag.note":     "Write the flag down and post it in Teams — that way we can see how far you are.",

		"cleanup.title": "The golden cleanup rule",
		"cleanup.body": `A container name exists only <strong>once</strong>. Before you start a container
			with the same name again, you have to remove the old one. Otherwise you get
			<em>"The container name is already in use"</em> — that is not your mistake!`,

		"footer.progress": "Your progress is stored as a cookie in this browser — not in the container. A private window therefore starts over.",
		"footer.reset":    "Completely tangled up? This clears everything and you restart at level 1:",
		"footer.resetcmd": "docker rm -f ctf-main data-provider-svc\ndocker network rm ctf-net",

		// ---- Level 1
		"l1.title": "The Secret Port",
		"l1.done": `You found the hidden port and published it. That is the only reason you can see
			this page at all. <strong>Level 1 complete!</strong>`,
		"l1.learned": `You can now look inside a container with <code>docker inspect</code> and publish
			a port with <code>-p</code>.`,

		// ---- Level 2
		"l2.title": "The Invisible Friend",
		"l2.goal": `I need a second container, the <strong>data provider</strong>. But I cannot reach
			it. Find out which <strong>name</strong> I am looking for, and start it with exactly
			that name.`,
		"l2.concept": `An image can carry <strong>labels</strong>: little notes the developers wrote into
			it. One of those notes is stuck to me and reveals what my friend must be called. You can
			see labels with <code>docker inspect</code>.`,
		"l2.hint1": `I carry metadata you can look at — just like you found the port in level 1. This
			time you are not looking for a port but for a <strong>label</strong>. It is called
			<code>ctf.data-provider.host</code>.`,
		"l2.hint2": `<code>docker inspect</code> shows you everything about a container — but that is
			easily 200 lines of JSON. With <code>--format</code> you pull out one specific value:
			<br><code>docker inspect --format '{{.Config.Labels}}' YOUR_CONTAINER_NAME</code>`,
		"l2.hint3": `Run this (your container is probably called <code>ctf-main</code>):
			<pre><code>docker inspect --format '{{.Config.Labels}}' ctf-main</code></pre>
			The output contains <code>ctf.data-provider.host:</code> followed by the name you need.
			Enter it below.`,
		"l2.check":       `Enter the correct name below and you get the level 2 flag.`,
		"l2.question":    "What must my friend, the data provider, be called?",
		"l2.placeholder": "Enter the name from the label …",
		"l2.submit":      "Check answer",
		"l2.wrong": `That is not the right name yet. Look at the labels —
			the key is called <code>ctf.data-provider.host</code>.`,
		"l2.flagtext": "You found the label and now know what the data provider must be called.",

		// ---- Level 3
		"l3.title": "The Network",
		"l3.goal": `The data provider has to be running <strong>and</strong> I have to be able to reach
			it. Make sure we are both on the same Docker network.
			<br><br>You will find its image in the Docker registry as
			<code>matixmedia/docker-ctf-data-provider:latest</code>.`,
		"l3.concept": `By default Docker puts all containers on one shared network — but
			<strong>there is no name resolution there</strong>. So containers cannot find each other by
			name. Only on a <strong>self-created</strong> network does a container's <code>--name</code>
			become its hostname.`,
		"l3.hint1": `The data provider lives in the image
			<code>matixmedia/docker-ctf-data-provider:latest</code>. It has to run, carry the right
			name — and we both have to be on the same network.`,
		"l3.hint2": `You need three building blocks:
			<br>1. Create a network: <code>docker network create NETNAME</code>
			<br>2. Start the data provider on it: <code>--name data-provider-svc --network NETNAME</code>
			<br>3. <strong>Restart me</strong>, also with <code>--network NETNAME</code>.
			<br>Step 3 is easy to forget — one container alone on the network is not enough!`,
		"l3.hint3": `Copy this whole block. The first two lines clean up in case something is already
			running — you can ignore any errors they produce.
			<pre><code>docker stop ctf-main data-provider-svc
docker rm ctf-main data-provider-svc

docker network create ctf-net

docker run -d --name data-provider-svc --network ctf-net \
  matixmedia/docker-ctf-data-provider:latest

docker run -d --name ctf-main --network ctf-net -p 8989:8989 \
  matixmedia/docker-ctf:latest</code></pre>
			Then reload this page.`,
		"l3.check": `As soon as it works, the red box above turns green and you get the flag.`,
		"l3.diag.dns": `<strong>I do not know that name.</strong> The name
			<code>data-provider-svc</code> cannot be resolved. That means we are
			<strong>not on the same network</strong>. That is exactly what this level is about.`,
		"l3.diag.refused": `<strong>I find the name, but nobody answers.</strong> The name resolves — so
			the container is not running, has crashed, or listens on a different port. Check with
			<code>docker ps</code> whether <code>data-provider-svc</code> is really running.`,
		"l3.flagtext": "Both containers are on the same network and find each other by name.",

		// ---- Level 4
		"l4.title": "The Vault",
		"l4.goal": `I know a password, but I only write it into a folder <em>inside</em> my container:
			<code>/secrets/password.txt</code>. Get it onto your machine and enter it below.`,
		"l4.concept": `A container is a closed box — delete it and every file inside is gone. With a
			<strong>bind mount</strong> (<code>-v</code>) you connect a folder on your machine to a
			folder in the container. Both then see <strong>the same</strong> files.`,
		"l4.hint1": `I only write the file if the folder <code>/secrets</code> actually exists. So you
			have to hand it to me from outside — and then you have to <strong>restart</strong> me,
			because I write the file at startup.`,
		"l4.hint2": `The option is <code>-v YOUR_FOLDER:/secrets</code>.
			Docker wants a <strong>full</strong> path here — <code>./secrets</code> is not enough.
			<code>$(pwd)</code> automatically inserts the folder you are currently in.`,
		"l4.hint3": `Copy this whole block:
			<pre><code>docker stop ctf-main
docker rm ctf-main

docker run -d --name ctf-main --network ctf-net -p 8989:8989 \
  -v $(pwd)/secrets:/secrets \
  matixmedia/docker-ctf:latest

cat secrets/password.txt</code></pre>
			The last command shows you the password. Enter it below.`,
		"l4.check": `After the restart you have a folder <code>secrets</code> containing
			<code>password.txt</code>.`,
		"l4.mounted.no":  `❌ I still do not see a <code>/secrets</code> folder. The mount is missing.`,
		"l4.mounted.yes": `✅ The folder <code>/secrets</code> is there — I wrote the password into it. Read it and enter it below.`,
		"l4.pwlabel":     "Password",
		"l4.placeholder": "Password from password.txt …",
		"l4.submit":      "Check password",
		"l4.wrong": `That password is wrong. Read the file again carefully:
			<code>cat secrets/password.txt</code> — no spaces before or after.`,
		"l4.flagtext": "You mounted a folder into the container and recovered the password.",

		// ---- Level 5
		"l5.title": "The Flag",
		"l5.goal": `I only show the final flag directly inside my container — and only if you bring a
			<strong>real terminal</strong>.`,
		"l5.concept": `With <code>docker exec</code> you run a command <em>inside</em> an already running
			container. The <code>-it</code> means <em>interactive</em> + <em>TTY</em>: it gives the
			command a real keyboard — without it this will not work.`,
		"l5.hint1": `You do not have to build me a new container. I am already running. You only have to
			execute a command <strong>inside me</strong>: <code>/app/app --show-flag</code>`,
		"l5.hint2": `The command starts with <code>docker exec ...</code> — but without
			<code>-it</code> I refuse to print anything. Feel free to try what happens without it.`,
		"l5.hint3": `The complete command:
			<pre><code>docker exec -it ctf-main /app/app --show-flag</code></pre>
			Then press <strong>ENTER</strong> once — and the flag is yours. 🎉`,
		"l5.check": `The final <code>FLAG{…}</code> appears in your terminal.`,
		"l5.bonus": `<strong>Bonus:</strong> feel free to look around inside me:
			<code>docker exec -it ctf-main sh</code> — with <code>ls</code> and <code>cat</code> you can
			walk through my file system. <code>exit</code> gets you back out.`,

		// ---- Done
		"done.title": "All levels complete!",
		"done.body": `You inspected images, published ports, named containers, built a network, mounted
			a volume and executed a command inside a running container. That is exactly what working
			with Docker looks like day to day. 🐳`,
	},
}
