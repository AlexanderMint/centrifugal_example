package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

const (
	jwtSecretKey = "a75f23a7-e9aa-4fd2-ab14-8825ff48bd13"
	port         = 8181
)

func main() {
	http.HandleFunc("/auth/anonymous", handleAnonymous)
	http.HandleFunc("/auth/anonymous/refresh", handleAnonymousRefresh)
	http.HandleFunc("/auth/subscribe", handleSubscribe)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handleAnonymous(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if !handleCORS(w, r) {
		return
	}

	user := uuid.New().String()
	token := connToken(user, time.Now().Add(time.Minute).Unix())
	json.NewEncoder(w).Encode(map[string]string{"token": token, "user_id": user})
}

func handleAnonymousRefresh(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if !handleCORS(w, r) {
		return
	}

	var req struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token := connToken(req.UserID, time.Now().Add(time.Minute).Unix())
	json.NewEncoder(w).Encode(map[string]string{"token": token, "user_id": req.UserID})
}

func handleSubscribe(w http.ResponseWriter, r *http.Request) {
	logRequest(r)
	if !handleCORS(w, r) {
		return
	}

	var req struct {
		UserID  string `json:"user_id"`
		Channel string `json:"channel"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token := subscriptionToken(req.Channel, req.UserID, time.Now().Add(time.Minute).Unix())
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func connToken(user string, exp int64) string {
	claims := jwt.MapClaims{"sub": user}
	if exp > 0 {
		claims["exp"] = exp
	}
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSecretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func subscriptionToken(channel string, user string, exp int64) string {
	claims := jwt.MapClaims{"channel": channel, "sub": user}
	if exp > 0 {
		claims["exp"] = exp
	}
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(jwtSecretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func logRequest(r *http.Request) {
	log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
}

func handleCORS(w http.ResponseWriter, r *http.Request) bool {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")

	return r.Method != http.MethodOptions
}
