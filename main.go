package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const defaultGreenAPI = "https://api.green-api.com"

func main() {
	greenBase := os.Getenv("GREEN_API_BASE")
	if greenBase == "" {
		greenBase = defaultGreenAPI
	}
	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/api/getSettings", makeHandler(greenBase, "getSettings"))
	http.HandleFunc("/api/getStateInstance", makeHandler(greenBase, "getStateInstance"))
	http.HandleFunc("/api/sendMessage", makeHandler(greenBase, "sendMessage"))
	http.HandleFunc("/api/sendFileByUrl", makeHandler(greenBase, "sendFileByUrl"))

	log.Printf("Server started on %s (GREEN_API_BASE=%s)\n", addr, greenBase)
	srv := &http.Server{Addr: addr, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second}
	log.Fatal(srv.ListenAndServe())
}

func makeHandler(greenBase, action string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		var payload map[string]interface{}
		if r.Body != nil {
			defer r.Body.Close()
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				payload = map[string]interface{}{}
			}
		} else {
			payload = map[string]interface{}{}
		}

		idI, _ := payload["idInstance"].(string)
		token, _ := payload["apiToken"].(string)
		if idI == "" || token == "" {
			http.Error(w, "idInstance and apiToken required in JSON body", http.StatusBadRequest)
			return
		}

		url := fmt.Sprintf("%s/waInstance%s/%s/%s", greenBase, idI, action, token)

		var resp *http.Response
		var err error
		if action == "getSettings" || action == "getStateInstance" {
			resp, err = http.Get(url)
		} else {
			delete(payload, "idInstance")
			delete(payload, "apiToken")
			bodyBytes, _ := json.Marshal(payload)
			req, rerr := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
			if rerr != nil {
				http.Error(w, "failed to create request: "+rerr.Error(), http.StatusInternalServerError)
				return
			}
			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{Timeout: 20 * time.Second}
			resp, err = client.Do(req)
		}

		if err != nil {
			http.Error(w, "request to GREEN-API failed: "+err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	}
}
