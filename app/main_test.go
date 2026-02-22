package main

import (
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "path/filepath"
    "testing"
)

func TestHandleRoot_ServesIndex(t *testing.T) {
    td := t.TempDir()
    // create static dir and index.html
    indexPath := filepath.Join(td, "index.html")
    content := []byte("<html><body>hello</body></html>")
    if err := ioutil.WriteFile(indexPath, content, 0644); err != nil {
        t.Fatal(err)
    }

    // point StaticDir to temp dir
    prev := StaticDir
    StaticDir = td
    defer func() { StaticDir = prev }()

    req := httptest.NewRequest("GET", "/", nil)
    w := httptest.NewRecorder()
    handleRoot(w, req)

    resp := w.Result()
    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 200 {
        t.Fatalf("expected 200, got %d", resp.StatusCode)
    }
    if string(body) != string(content) {
        t.Fatalf("unexpected body: %s", string(body))
    }
}

func TestHandleRoot_ServesStaticFile(t *testing.T) {
    td := t.TempDir()
    // create static dir and a css file
    cssPath := filepath.Join(td, "style.css")
    content := []byte("body { color: red }")
    if err := ioutil.WriteFile(cssPath, content, 0644); err != nil {
        t.Fatal(err)
    }

    prev := StaticDir
    StaticDir = td
    defer func() { StaticDir = prev }()

    req := httptest.NewRequest("GET", "/style.css", nil)
    w := httptest.NewRecorder()
    handleRoot(w, req)

    resp := w.Result()
    body, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != 200 {
        t.Fatalf("expected 200, got %d", resp.StatusCode)
    }
    if string(body) != string(content) {
        t.Fatalf("unexpected body: %s", string(body))
    }
}

func TestHandleRoot_RedirectsWhenNotFoundAndNoClient(t *testing.T) {
    td := t.TempDir()
    prev := StaticDir
    StaticDir = td
    defer func() { StaticDir = prev }()

    // ensure no store is configured
    prevStore := store
    store = nil
    defer func() { store = prevStore }()

    req := httptest.NewRequest("GET", "/not-exists", nil)
    w := httptest.NewRecorder()
    handleRoot(w, req)

    resp := w.Result()
    if resp.StatusCode != http.StatusSeeOther {
        t.Fatalf("expected redirect (303), got %d", resp.StatusCode)
    }
    loc, _ := resp.Location()
    if loc.Path != "/" {
        t.Fatalf("expected redirect to /, got %s", loc.Path)
    }
}
