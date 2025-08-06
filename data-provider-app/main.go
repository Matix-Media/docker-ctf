package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	correctPassword = "SUPER_GEHEIM_123"
	// Aktualisierter Hinweis für den letzten Schritt.
	finalHint       = "FAST GESCHAFFT: Um die Flagge zu sehen, musst du den Hauptcontainer anweisen, sie direkt auszugeben. Führe im Hauptcontainer den Befehl '/app/app --show-flag' aus."
)

func main() {
	// Ein einfacher Ping-Endpunkt, um die Erreichbarkeit zu testen
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "pong")
	})

	// Ein Endpunkt, der das Passwort überprüft
	http.HandleFunc("/verify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Nur POST-Anfragen erlaubt", http.StatusMethodNotAllowed)
			return
		}

		submittedPassword := r.FormValue("password")
		if submittedPassword == correctPassword {
			fmt.Fprint(w, finalHint)
		} else {
			http.Error(w, "Falsches Passwort!", http.StatusUnauthorized)
		}
	})

	log.Println("Data-Provider Service startet auf Port 9090...")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatalf("Konnte den Data-Provider nicht starten: %s\n", err)
	}
}
