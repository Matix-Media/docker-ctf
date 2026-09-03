package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const correctPassword = "SUPER_GEHEIM_123"

// Antworten in beiden Sprachen. Die Sprache kommt von der Main-App mit.
var messages = map[string]map[string]string{
	"de": {
		"hint": "FAST GESCHAFFT: Fuehre im Hauptcontainer 'docker exec -it ctf-main /app/app --show-flag' aus.",
		"bad":  "Falsches Passwort! Schau nochmal in secrets/password.txt.",
	},
	"en": {
		"hint": "ALMOST THERE: run 'docker exec -it ctf-main /app/app --show-flag' on the main container.",
		"bad":  "Wrong password! Check secrets/password.txt again.",
	},
}

func msg(lang, key string) string {
	if m, ok := messages[lang]; ok {
		if v, ok := m[key]; ok {
			return v
		}
	}
	return messages["de"][key]
}

func main() {
	// Erreichbarkeitstest fuer die Main-App (Level 3).
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "pong")
	})

	// Passwortpruefung (Level 4).
	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Nur POST-Anfragen erlaubt / POST only", http.StatusMethodNotAllowed)
			return
		}

		lang := r.FormValue("lang")
		if lang != "en" {
			lang = "de"
		}

		if strings.TrimSpace(r.FormValue("password")) == correctPassword {
			fmt.Fprint(w, msg(lang, "hint"))
			return
		}
		http.Error(w, msg(lang, "bad"), http.StatusUnauthorized)
	})

	log.Println("Data-Provider startet auf Port 9090 / starting on port 9090...")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatalf("Konnte den Data-Provider nicht starten / could not start: %s\n", err)
	}
}
