package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	// Die Adresse des zweiten Containers.
	dataProviderHost = "data-provider-svc:9090"
	// Das Passwort, das in das Volume geschrieben wird.
	password = "SUPER_GEHEIM_123"
)

// main ist der Einstiegspunkt des Programms.
func main() {
	// Schritt 8: Füge eine Kommandozeilen-Flag hinzu, um die finale Flagge anzuzeigen.
	showFlag := flag.Bool("show-flag", false, "Zeigt die finale Flagge an und beendet das Programm.")
	flag.Parse()

	if *showFlag {
		// Prüfe, ob der Prozess in einem interaktiven Terminal läuft.
		stat, err := os.Stdin.Stat()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fehler: Konnte den Status des Terminals nicht ermitteln (Grund: %v).\n", err)
			os.Exit(1)
		}

		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// Kein Terminal: Gib eine Fehlermeldung aus und beende.
			fmt.Fprintln(os.Stderr, "Fehler: Dies muss in einem interaktiven Terminal ausgeführt werden.")
			fmt.Fprintln(os.Stderr, "Benutze 'docker exec -it ...'.")
			os.Exit(1)
		}

		// Terminal ist vorhanden: Warte auf die Eingabe.
		fmt.Println("Du bist fast am Ziel! Drücke ENTER, um die Flagge anzuzeigen.")
		
		// Lese die Eingabe. Wenn der Input-Kanal geschlossen ist (wie bei `docker exec` ohne -i),
		// gibt dieser Befehl sofort einen Fehler zurück.
		_, err = bufio.NewReader(os.Stdin).ReadByte()
		if err != nil {
			// Wenn ein Fehler auftritt (z.B. EOF), gib eine spezifische Fehlermeldung aus.
			fmt.Fprintf(os.Stderr, "\nFehler: Konnte keine Eingabe vom Terminal lesen (Grund: %v).\n", err)
			fmt.Fprintln(os.Stderr, "Stelle sicher, dass du eine interaktive Sitzung mit '-it' gestartet hast.")
			os.Exit(1)
		}
		fmt.Println("FLAG{D0CK3R_PR0F1_MIT_FLAG}")
		return
	}

	// Schritt 1: Ein Hinweis, der zu `docker inspect` führt.
	log.Println("HINWEIS: Ich lausche auf einem geheimen Port. Finde ihn mit 'docker container inspect DEIN_CONTAINER' und schau unter 'Config.ExposedPorts'.")

	// Schritt 6: Schreibe das Passwort in das Volume, falls es gemountet ist.
	writePasswordToVolume()

	// Definiere den Handler für die Hauptroute "/".
	http.HandleFunc("/", rootHandler)

	// Starte den Webserver auf einem nicht-standard Port.
	fmt.Println("Server startet auf Port 8989...")
	if err := http.ListenAndServe(":8989", nil); err != nil {
		log.Fatalf("Konnte den Server nicht starten: %s\n", err)
	}
}

// writePasswordToVolume prüft, ob das Verzeichnis /secrets existiert und schreibt die Passwortdatei.
func writePasswordToVolume() {
	secretsDir := "/secrets"
	if _, err := os.Stat(secretsDir); !os.IsNotExist(err) {
		filePath := filepath.Join(secretsDir, "password.txt")
		err := os.WriteFile(filePath, []byte(password), 0644)
		if err != nil {
			log.Printf("Konnte Passwort nicht in Volume schreiben: %v", err)
		} else {
			log.Printf("Passwort erfolgreich nach %s geschrieben.", filePath)
		}
	}
}

// rootHandler behandelt Anfragen und die Passworteingabe.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	// Wenn ein Passwort gesendet wird, verarbeite es.
	if r.Method == http.MethodPost {
		handlePasswordSubmission(w, r)
		return
	}

	// Schritt 4 & 5: Prüfe die Verbindung zum Data-Provider.
	resp, err := http.Get("http://" + dataProviderHost + "/ping")
	connectionOK := err == nil && resp.StatusCode == http.StatusOK
	if err == nil {
		resp.Body.Close()
	}

	pageData := struct {
		ConnectionOK bool
		Message      string
	}{
		ConnectionOK: connectionOK,
		Message:      "",
	}

	// Lade die Webseite.
	tmpl, err := template.New("index").Parse(getIndexTemplate())
	if err != nil {
		http.Error(w, "Konnte Template nicht parsen", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, pageData)
}

// handlePasswordSubmission verarbeitet das gesendete Passwort.
func handlePasswordSubmission(w http.ResponseWriter, r *http.Request) {
	submittedPassword := r.FormValue("password")

	// Schritt 7: Sende das Passwort an den Data-Provider zur Verifizierung.
	apiURL := fmt.Sprintf("http://%s/verify", dataProviderHost)
	postBody := []byte(fmt.Sprintf("password=%s", submittedPassword))
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(postBody))
	if err != nil {
		http.Error(w, "Fehler beim Erstellen der Anfrage", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Fehler bei der Verbindung zum Data-Provider", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	// Zeige die Antwort vom Data-Provider an.
	pageData := struct {
		ConnectionOK bool
		Message      string
	}{
		ConnectionOK: true,
		Message:      string(body),
	}

	tmpl, err := template.New("index").Parse(getIndexTemplate())
	if err != nil {
		http.Error(w, "Konnte Template nicht parsen", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tmpl.Execute(w, pageData)
}

// getIndexTemplate gibt das HTML-Template für die Webseite zurück.
func getIndexTemplate() string {
	return `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Docker CTF</title>
		<style>
			body { font-family: sans-serif; background-color: #f0f0f0; color: #333; max-width: 800px; margin: 40px auto; padding: 20px; border-radius: 8px; box-shadow: 0 4px 8px rgba(0,0,0,0.1); }
			h1 { color: #005a9e; }
			code { background-color: #e0e0e0; padding: 2px 6px; border-radius: 4px; }
			.hint { border-left: 4px solid #ffc107; padding: 10px; margin-top: 20px; }
			.success { border-left: 4px solid #28a745; padding: 10px; margin-top: 20px; margin-bottom: 20px; }
			.error { border-left: 4px solid #dc3545; padding: 10px; margin-top: 20px; }
			input[type=text], button { padding: 10px; margin-top: 10px; border-radius: 4px; border: 1px solid #ccc; width: calc(100% - 22px); }
			button { background-color: #007bff; color: white; cursor: pointer; width: 100%; }
		</style>
	</head>
	<body>
		<h1>Willkommen beim Docker CTF!</h1>
		
		{{if .ConnectionOK}}
			<div class="success">
				<strong>Verbindung zum Data-Provider erfolgreich!</strong>
				<p>Jetzt brauchst du das richtige Passwort, um den letzten Hinweis zu erhalten.</p>
				<p><strong>Hinweis:</strong> Ich habe das Passwort an einen sicheren Ort geschrieben: <code>/secrets/password.txt</code>. Du musst diesen Ordner nur von deinem eigenen Computer aus zugänglich machen. Starte den Container neu und benutze die <code>-v</code> Option, um ein Volume zu mounten. Z.B.: <code>-v $(pwd)/secrets:/secrets</code></p>
			</div>

			<form method="POST" action="/">
				<label for="password">Passwort:</label><br>
				<input type="text" id="password" name="password"><br>
				<button type="submit">Absenden</button>
			</form>

			{{if .Message}}
				<div class="hint">
					<h2>Antwort vom Data-Provider:</h2>
					<p><code>{{.Message}}</code></p>
				</div>
			{{end}}

		{{else}}
			<div class="error">
				<strong>Verbindung zum Data-Provider fehlgeschlagen!</strong>
				<p>Ich muss mit meinem Freund, dem Data-Provider, kommunizieren, aber ich kann ihn nicht erreichen. Das Image für den Data-Provider ist in der Docker Registry unter dem Namen <code>matixmedia/docker-ctf-data-provider:latest</code> zu finden.</p>
				<p><strong>Hinweis:</strong> Container im Standard-Netzwerk können sich nicht über ihre Namen erreichen. Du musst ein eigenes Docker-Netzwerk erstellen (<code>docker network create ctf-net</code>) und beide Container in diesem Netzwerk starten (<code>--network ctf-net</code>).</p>
				<p>Den Hostnamen des Data-Providers findest du übrigens als Label in meinen Metadaten. Benutze <code>docker container inspect DEIN_CONTAINER_NAME</code> und suche nach <code>ctf.data-provider.host</code>.</p>
			</div>
		{{end}}
	</body>
	</html>
	`
}
