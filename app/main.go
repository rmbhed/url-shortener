package main

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strings"

    "cloud.google.com/go/firestore"
)

var store LinkStore
// StaticDir is the directory that holds static frontend files. Made configurable for tests.
var StaticDir = "./static"

func main() {
    ctx := context.Background()
    projectID := os.Getenv("PROJECT_ID")

    // Initialize Firestore-backed store
    client, err := firestore.NewClient(ctx, projectID)
    if err != nil {
        log.Fatalf("Failed to create Firestore client: %v", err)
    }
    defer client.Close()

    store = &FirestoreStore{Client: client}

    // API Handlers
    http.HandleFunc("/api/links", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodGet {
            handleListLinks(w, r)
        } else if r.Method == http.MethodPost {
            handleCreateLink(w, r)
        }
    })

    // Everything else (UI + Redirects)
    http.HandleFunc("/", handleRoot)

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Listening on port %s", port)
    if err := http.ListenAndServe(":"+port, nil); err != nil {
        log.Fatal(err)
    }
}

func handleListLinks(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    if store == nil {
        http.Error(w, "store not configured", http.StatusInternalServerError)
        return
    }
    links, err := store.GetAll(ctx)
    if err != nil {
        http.Error(w, "Failed to fetch links", http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(links)
}

func handleCreateLink(w http.ResponseWriter, r *http.Request) {
    var l Link
    if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
    if store == nil {
        http.Error(w, "store not configured", http.StatusInternalServerError)
        return
    }
    if err := store.Create(r.Context(), l); err != nil {
        // Conflict (duplicate)
        w.WriteHeader(http.StatusConflict)
        return
    }
    w.WriteHeader(http.StatusCreated)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
    // 1. If it's just the root path "/", serve the frontend
    if r.URL.Path == "/" {
        http.ServeFile(w, r, filepath.Join(StaticDir, "index.html"))
        return
    }

    // 2. Check if it's a file that exists (like style.css or script.js)
    path := filepath.Join(StaticDir, r.URL.Path)
    if _, err := os.Stat(path); err == nil {
        http.ServeFile(w, r, path)
        return
    }

    // 3. Assume it's a short link: Check store
    shortName := strings.TrimPrefix(r.URL.Path, "/")
    if store == nil {
        // In tests or when store isn't configured, redirect to home
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    l, err := store.Get(r.Context(), shortName)
    if err != nil {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }
    http.Redirect(w, r, l.URL, http.StatusMovedPermanently)
}