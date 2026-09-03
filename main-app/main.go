package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	// Adresse des zweiten Containers. / Address of the second container.
	dataProviderHost = "data-provider-svc:9090"
	// Passwort, das in das Volume geschrieben wird. / Password written into the volume.
	password = "SUPER_GEHEIM_123"
	// Wert des Labels ctf.data-provider.host — die Antwort für Level 2.
	providerLabelValue = "data-provider-svc"
	// Ordner, der in Level 4 gemountet werden muss.
	secretsDir = "/secrets"

	cookieLang     = "ctf_lang"
	cookieProgress = "ctf_progress"
)

// Zwischen-Flags pro Level. Level 5 gibt es nur im Terminal.
var levelFlags = map[int]string{
	1: "FLAG{L1_P0RT_G3FUNDEN}",
	2: "FLAG{L2_L4B3L_G3L3S3N}",
	3: "FLAG{L3_N3TZW3RK_ST3HT}",
	4: "FLAG{L4_TR3S0R_G3KN4CKT}",
}

const finalFlag = "FLAG{D0CK3R_PR0F1_MIT_FLAG}"

// providerStatus beschreibt, warum der Data-Provider (nicht) erreichbar ist.
type providerStatus int

const (
	providerOK      providerStatus = iota // Verbindung steht
	providerDNS                           // Name nicht auflösbar -> Netzwerkproblem
	providerRefused                       // Name bekannt, aber niemand antwortet
)

func main() {
	showFlag := flag.Bool("show-flag", false, "Zeigt die finale Flagge an / shows the final flag")
	flag.Parse()

	if *showFlag {
		runShowFlag()
		return
	}

	printStartupHints()
	writePasswordToVolume()

	http.HandleFunc("/", rootHandler)

	// Bewusst ohne Portnummer: die soll in Level 1 selbst gefunden werden.
	fmt.Println("Server laeuft. / Server is running.")
	if err := http.ListenAndServe(":8989", nil); err != nil {
		log.Fatalf("Konnte den Server nicht starten / could not start server: %s\n", err)
	}
}

// runShowFlag ist Level 5: die Flagge gibt es nur in einem echten Terminal.
func runShowFlag() {
	if !stdinIsTerminal() {
		printNoTTYHelp()
		os.Exit(1)
	}

	fmt.Println("Du bist fast am Ziel! Drücke ENTER, um die Flagge anzuzeigen.")
	fmt.Println("You are almost there! Press ENTER to reveal the flag.")

	if _, err := bufio.NewReader(os.Stdin).ReadByte(); err != nil {
		printNoTTYHelp()
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("  " + finalFlag)
	fmt.Println()
	fmt.Println("Glückwunsch, du hast alle Level geschafft! 🐳")
	fmt.Println("Congratulations, you completed every level! 🐳")
}

// printNoTTYHelp erklärt, was fehlt — und nennt den kompletten richtigen Befehl.
func printNoTTYHelp() {
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  Fast! Mir fehlt ein echtes Terminal.")
	fmt.Fprintln(os.Stderr, "  Almost! I am missing a real terminal.")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  Benutze genau diesen Befehl / use exactly this command:")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "      docker exec -it ctf-main /app/app --show-flag")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "  Das -it ist der Unterschied: i = interactive, t = terminal.")
	fmt.Fprintln(os.Stderr, "  The -it is what matters: i = interactive, t = terminal.")
	fmt.Fprintln(os.Stderr, "")
}

// printStartupHints ist der einzige Hinweis fuer Level 1.
//
// WICHTIG: Die Spieler bekommen nur den Image-Namen genannt - kein Repository,
// keine Anleitung, keine Skripte. Alles, was sie brauchen, muss also hier oder
// spaeter auf der Webseite stehen. Dieser Text ist der komplette Einstieg.
func printStartupHints() {
	lines := []string{
		"",
		"  ==================================================================",
		"   Docker CTF  -  Level 1 von 5:  Der geheime Port",
		"   Docker CTF  -  Level 1 of 5:   The Secret Port",
		"  ==================================================================",
		"",
		"  --- DEUTSCH ------------------------------------------------------",
		"",
		"  Willkommen! Du loest 5 Level und bekommst fuer jedes eine Flag.",
		"  Ab Level 2 fuehrt dich eine Webseite weiter - die erreichst du",
		"  aber erst, wenn du meinen Port gefunden und freigegeben hast.",
		"",
		"   1) Starte mich im Hintergrund und gib mir einen Namen:",
		"        docker run -d --name ctf-main matixmedia/docker-ctf",
		"",
		"      (Laeuft dieses Fenster gerade fest? Dann hast du mich ohne -d",
		"       gestartet. Druecke Strg+C und benutze den Befehl oben.)",
		"",
		"   2) Finde heraus, auf welchem Port ich lausche:",
		"        docker inspect --format '{{.Config.ExposedPorts}}' ctf-main",
		"",
		"   3) Starte mich neu und gib den Port frei.",
		"      PORT ist die Zahl aus Schritt 2:",
		"        docker stop ctf-main && docker rm ctf-main",
		"        docker run -d --name ctf-main -p PORT:PORT matixmedia/docker-ctf",
		"",
		"   4) Oeffne im Browser:  http://localhost:PORT",
		"",
		"  MERKE: Einen Containernamen gibt es nur einmal. Vor jedem Neustart",
		"  mit gleichem Namen erst 'docker stop NAME && docker rm NAME'.",
		"  Weisst du meinen Namen nicht mehr? 'docker ps' zeigt ihn dir.",
		"",
		"  --- ENGLISH ------------------------------------------------------",
		"",
		"  Welcome! You solve 5 levels and get a flag for each one.",
		"  From level 2 on, a web page guides you - but you can only reach it",
		"  once you have found and published my port.",
		"",
		"   1) Start me in the background and give me a name:",
		"        docker run -d --name ctf-main matixmedia/docker-ctf",
		"",
		"      (Is this window stuck? Then you started me without -d.",
		"       Press Ctrl+C and use the command above.)",
		"",
		"   2) Find out which port I am listening on:",
		"        docker inspect --format '{{.Config.ExposedPorts}}' ctf-main",
		"",
		"   3) Restart me and publish the port.",
		"      PORT is the number from step 2:",
		"        docker stop ctf-main && docker rm ctf-main",
		"        docker run -d --name ctf-main -p PORT:PORT matixmedia/docker-ctf",
		"",
		"   4) Open in your browser:  http://localhost:PORT",
		"",
		"  REMEMBER: a container name exists only once. Before restarting with",
		"  the same name, run 'docker stop NAME && docker rm NAME' first.",
		"  Forgot my name? 'docker ps' shows it.",
		"",
		"  ==================================================================",
		"",
	}
	for _, l := range lines {
		fmt.Println(l)
	}
}

// writePasswordToVolume schreibt das Passwort, sobald /secrets existiert (Level 4).
func writePasswordToVolume() {
	if _, err := os.Stat(secretsDir); os.IsNotExist(err) {
		return
	}
	filePath := filepath.Join(secretsDir, "password.txt")
	if err := os.WriteFile(filePath, []byte(password), 0644); err != nil {
		log.Printf("Konnte Passwort nicht in Volume schreiben / could not write password: %v", err)
		return
	}
	log.Printf("Passwort nach %s geschrieben / password written to %s", filePath, filePath)
}

// ---------------------------------------------------------------- Zustand

// checkProvider prueft die Verbindung und unterscheidet die Fehlerursache.
func checkProvider() providerStatus {
	resp, err := http.Get("http://" + dataProviderHost + "/ping")
	if err == nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			return providerOK
		}
		return providerRefused
	}

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) || strings.Contains(err.Error(), "no such host") {
		return providerDNS
	}
	return providerRefused
}

func volumeMounted() bool {
	_, err := os.Stat(secretsDir)
	return err == nil
}

func langFrom(r *http.Request) Lang {
	if q := r.URL.Query().Get("lang"); q != "" {
		return ParseLang(q)
	}
	if c, err := r.Cookie(cookieLang); err == nil {
		return ParseLang(c.Value)
	}
	return LangDE
}

func progressFrom(r *http.Request) int {
	c, err := r.Cookie(cookieProgress)
	if err != nil {
		return 0
	}
	n, err := strconv.Atoi(c.Value)
	if err != nil || n < 0 || n > 5 {
		return 0
	}
	return n
}

func setCookie(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:   name,
		Value:  value,
		Path:   "/",
		MaxAge: 60 * 60 * 24 * 30,
	})
}

// ---------------------------------------------------------------- Handler

func rootHandler(w http.ResponseWriter, r *http.Request) {
	lang := langFrom(r)
	setCookie(w, cookieLang, string(lang))

	stored := progressFrom(r)
	progress := stored
	status := checkProvider()

	// Level 3 ist automatisch geschafft, sobald die Verbindung steht.
	if status == providerOK && progress < 3 {
		progress = 3
	}
	// Wer diese Seite sieht, hat Level 1 hinter sich.
	if progress < 1 {
		progress = 1
	}

	var justUnlocked int
	var wrongKey string

	if r.Method == http.MethodPost {
		switch r.FormValue("step") {
		case "2":
			answer := strings.TrimSpace(r.FormValue("answer"))
			if strings.EqualFold(answer, providerLabelValue) {
				if progress < 2 {
					progress = 2
					justUnlocked = 2
				}
			} else if answer != "" {
				wrongKey = "l2.wrong"
			}
		case "4":
			submitted := strings.TrimSpace(r.FormValue("password"))
			ok, err := verifyPassword(submitted, lang)
			switch {
			case err != nil:
				wrongKey = "l3.diag.refused"
			case ok:
				if progress < 4 {
					progress = 4
					justUnlocked = 4
				}
			case submitted != "":
				wrongKey = "l4.wrong"
			}
		}
	}

	// Level 3 frisch geschafft? Nur beim ersten Mal hervorheben.
	if justUnlocked == 0 && status == providerOK && stored < 3 {
		justUnlocked = 3
	}

	setCookie(w, cookieProgress, strconv.Itoa(progress))
	render(w, lang, progress, status, justUnlocked, wrongKey)
}

// verifyPassword laesst das Passwort vom Data-Provider pruefen (Level 4).
func verifyPassword(submitted string, lang Lang) (bool, error) {
	form := url.Values{}
	form.Set("password", submitted)
	form.Set("lang", string(lang))

	apiURL := "http://" + dataProviderHost + "/verify"
	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBufferString(form.Encode()))
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	return resp.StatusCode == http.StatusOK, nil
}

// ---------------------------------------------------------------- Rendering

type levelDot struct {
	Number int
	State  string // "done" | "current" | "todo"
}

type earnedFlag struct {
	Level int
	Title string
	Text  template.HTML
	Value string
	Fresh bool
}

type pageData struct {
	Lang       Lang
	OtherLang  Lang
	Dots       []levelDot
	Flags      []earnedFlag
	Connected  bool
	Diagnosis  template.HTML
	Level      int
	LevelTitle string
	Goal       template.HTML
	Concept    template.HTML
	Hint1      template.HTML
	Hint2      template.HTML
	Hint3      template.HTML
	Check      template.HTML
	Bonus      template.HTML
	ShowForm2  bool
	ShowForm4  bool
	Mounted    bool
	MountNote  template.HTML
	Error      template.HTML

	// Statische Beschriftungen
	Brand          string
	LangSwitch     string
	ProgressWord   string
	LevelWord      string
	GoalLabel      string
	ConceptLabel   string
	CheckLabel     string
	Hint1Label     string
	Hint2Label     string
	Hint3Label     string
	FlagUnlocked   string
	FlagNote       string
	CleanupTitle   string
	CleanupBody    template.HTML
	FooterProgress string
	FooterReset    string
	FooterResetCmd string
	DoneTitle      string
	DoneBody       template.HTML
	L1Done         template.HTML
	Question       string
	Placeholder2   string
	Submit2        string
	PwLabel        string
	Placeholder4   string
	Submit4        string
	HintsLabel     string
	ConnOK         string
	ConnBad        string
	ConnTitle      string
}

// currentLevel leitet aus Fortschritt und Verbindungszustand das anzuzeigende Level ab.
// Bricht die Verbindung wieder weg, landet man automatisch zurueck bei Level 3.
func currentLevel(progress int, status providerStatus) int {
	if status != providerOK {
		if progress < 2 {
			return 2
		}
		return 3
	}
	if progress < 4 {
		return 4
	}
	return 5
}

func render(w http.ResponseWriter, lang Lang, progress int, status providerStatus, justUnlocked int, wrongKey string) {
	level := currentLevel(progress, status)

	dots := make([]levelDot, 0, 5)
	for i := 1; i <= 5; i++ {
		state := "todo"
		switch {
		case i <= progress:
			state = "done"
		case i == level:
			state = "current"
		}
		dots = append(dots, levelDot{Number: i, State: state})
	}

	flags := make([]earnedFlag, 0, 4)
	for i := 1; i <= progress && i <= 4; i++ {
		textKey := fmt.Sprintf("l%d.flagtext", i)
		if i == 1 {
			textKey = "l1.learned"
		}
		flags = append(flags, earnedFlag{
			Level: i,
			Title: T(lang, fmt.Sprintf("l%d.title", i)),
			Text:  TH(lang, textKey),
			Value: levelFlags[i],
			Fresh: i == justUnlocked,
		})
	}

	data := pageData{
		Lang:      lang,
		OtherLang: lang.Other(),
		Dots:      dots,
		Flags:     flags,
		Connected: status == providerOK,
		Level:     level,
		Mounted:   volumeMounted(),

		Brand:          T(lang, "brand"),
		LangSwitch:     T(lang, "lang.switch"),
		ProgressWord:   T(lang, "progress.label"),
		LevelWord:      T(lang, "level.word"),
		GoalLabel:      T(lang, "goal.label"),
		ConceptLabel:   T(lang, "concept.label"),
		CheckLabel:     T(lang, "check.label"),
		Hint1Label:     T(lang, "hint.1"),
		Hint2Label:     T(lang, "hint.2"),
		Hint3Label:     T(lang, "hint.3"),
		FlagUnlocked:   T(lang, "flag.unlocked"),
		FlagNote:       T(lang, "flag.note"),
		CleanupTitle:   T(lang, "cleanup.title"),
		CleanupBody:    TH(lang, "cleanup.body"),
		FooterProgress: T(lang, "footer.progress"),
		FooterReset:    T(lang, "footer.reset"),
		FooterResetCmd: T(lang, "footer.resetcmd"),
		DoneTitle:      T(lang, "done.title"),
		DoneBody:       TH(lang, "done.body"),
		L1Done:         TH(lang, "l1.done"),
		Question:       T(lang, "l2.question"),
		Placeholder2:   T(lang, "l2.placeholder"),
		Submit2:        T(lang, "l2.submit"),
		PwLabel:        T(lang, "l4.pwlabel"),
		Placeholder4:   T(lang, "l4.placeholder"),
		Submit4:        T(lang, "l4.submit"),
		HintsLabel:     T(lang, "hints.label"),
		ConnOK:         T(lang, "conn.ok"),
		ConnBad:        T(lang, "conn.bad"),
	}

	// Der Hostname des Data-Providers IST die Loesung von Level 2. Vor Level 3
	// darf ihn also weder die Statusbox noch die Diagnose noch der Aufraeum-
	// Befehl im Footer verraten.
	revealProvider := level >= 3
	if revealProvider {
		data.ConnTitle = providerLabelValue
		switch status {
		case providerDNS:
			data.Diagnosis = TH(lang, "l3.diag.dns")
		case providerRefused:
			data.Diagnosis = TH(lang, "l3.diag.refused")
		}
	} else {
		data.ConnTitle = T(lang, "conn.anon.title")
		data.ConnBad = T(lang, "conn.anon.bad")
		data.FooterResetCmd = T(lang, "footer.resetcmd.l2")
	}

	prefix := fmt.Sprintf("l%d.", level)
	data.LevelTitle = T(lang, prefix+"title")
	data.Goal = TH(lang, prefix+"goal")
	data.Concept = TH(lang, prefix+"concept")
	data.Hint1 = TH(lang, prefix+"hint1")
	data.Hint2 = TH(lang, prefix+"hint2")
	data.Hint3 = TH(lang, prefix+"hint3")
	data.Check = TH(lang, prefix+"check")

	switch level {
	case 2:
		data.ShowForm2 = true
	case 4:
		data.ShowForm4 = true
		if data.Mounted {
			data.MountNote = TH(lang, "l4.mounted.yes")
		} else {
			data.MountNote = TH(lang, "l4.mounted.no")
		}
	case 5:
		data.Bonus = TH(lang, "l5.bonus")
	}

	if wrongKey != "" {
		data.Error = TH(lang, wrongKey)
	}

	tmpl, err := template.New("index").Parse(indexTemplate)
	if err != nil {
		http.Error(w, "template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("render error: %v", err)
	}
}
