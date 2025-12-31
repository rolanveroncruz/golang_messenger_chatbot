package main

/// ---- Main Webhook -----
/// Meta calls something like:
/// GET /webhook?hub.mode=subscribe&hub.verify_token=...&hub.challenge=...
/// if hub.verify_token matches my verify token, respond 200 with the raw hub.challenge string.
/// If not, respond 403.

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rolanveroncruz/golang_messenger_chatbot/aiChat"
)

func main() {
	_ = godotenv.Load()

	verifyToken := os.Getenv("FB_VERIFY_TOKEN")
	if verifyToken == "" {
		log.Fatal("FB_VERIFY_TOKEN is required (set it to the same value you type in Meta Webhooks Verify Token)")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := mux.NewRouter()

	// Meta webhook endpoint: same path handles GET (verification) + POST (events)
	r.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleVerify(w, r, verifyToken)
		case http.MethodPost:
			handleEvent(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}).Methods(http.MethodGet, http.MethodPost)

	// Optional: quick health check
	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}).Methods(http.MethodGet)

	addr := ":" + port
	log.Printf("Listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

// GET /webhook?hub.mode=subscribe&hub.verify_token=...&hub.challenge=...
func handleVerify(w http.ResponseWriter, r *http.Request, verifyToken string) {
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")
	fmt.Printf("%s ***handleVerify start***\n", time.Now().Format(time.RFC3339))
	println(mode)
	println(token)
	println(challenge)
	println("***handleVerify end***")

	if mode == "subscribe" && token == verifyToken {
		// IMPORTANT: Must respond with the raw challenge string (not JSON).
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(challenge))
		return
	}

	http.Error(w, "forbidden", http.StatusForbidden)
}

// POST /webhook (Messenger will send JSON events here)
func handleEvent(w http.ResponseWriter, r *http.Request) {
	// For now, just read and log the body so you can see what Meta sends.
	// Later you'll parse JSON and call the Send API to respond.
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(r.Body)

	log.Printf("Webhook event: %s\n", string(raw))

	var payload aiChat.WebhookPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("Failed to parsse webhook JSON: %v\n", err)
		//Always respond with 200 OK
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "EVENT_RECEIVED")

		// Meta sends batches: entry[] -> messaging[]
		for _, entry := range payload.Entry {
			for _, ev := range entry.Messaging {
				// 1. Handle text message
				if ev.Message != nil {
					// Ignore echoes of messages (prevents echo loops)
					if ev.Message.IsEcho {
						continue
					}
					if ev.Message.Text != "" {
						psid := ev.Sender.ID
						text := ev.Message.Text

						log.Printf("Got message from %s: %q\n", text, psid)
						// Reply with echo
						if err := aiChat.SendText(psid, "echo: "+text); err != nil {
							log.Printf("Failed to send message: %v\n", err)
						}
					}
					continue
				}
			}
		}
		// Must respond with 200 OK so Meta doesn't retry.
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, "EVENT_RECEIVED")
	}
}
