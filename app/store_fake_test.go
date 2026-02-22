package main

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "sync"
    "testing"
)

type FakeStore struct{
    mu sync.Mutex
    data map[string]Link
}

func NewFakeStore() *FakeStore {
    return &FakeStore{data: map[string]Link{}}
}

func (f *FakeStore) GetAll(ctx context.Context) ([]Link, error) {
    f.mu.Lock()
    defer f.mu.Unlock()
    out := make([]Link, 0, len(f.data))
    for _, v := range f.data { out = append(out, v) }
    return out, nil
}

func (f *FakeStore) Create(ctx context.Context, l Link) error {
    f.mu.Lock()
    defer f.mu.Unlock()
    if _, ok := f.data[l.ShortName]; ok { return &ConflictError{} }
    f.data[l.ShortName] = l
    return nil
}

func (f *FakeStore) Get(ctx context.Context, shortName string) (Link, error) {
    f.mu.Lock()
    defer f.mu.Unlock()
    v, ok := f.data[shortName]
    if !ok { return Link{}, ErrNotFound }
    return v, nil
}

// Simple sentinel errors for fake store
type ConflictError struct{}
func (ConflictError) Error() string { return "conflict" }
var ErrNotFound = &NotFoundError{}
type NotFoundError struct{}
func (NotFoundError) Error() string { return "not found" }

func TestCreateAndList(t *testing.T) {
    f := NewFakeStore()
    prev := store
    store = f
    defer func(){ store = prev }()

    // Create a link
    l := Link{ShortName: "foo", URL: "https://example.com"}
    body, _ := json.Marshal(l)
    req := httptest.NewRequest("POST", "/api/links", bytes.NewReader(body))
    w := httptest.NewRecorder()
    handleCreateLink(w, req)
    resp := w.Result()
    if resp.StatusCode != http.StatusCreated {
        t.Fatalf("expected 201 created, got %d", resp.StatusCode)
    }

    // List links
    req2 := httptest.NewRequest("GET", "/api/links", nil)
    w2 := httptest.NewRecorder()
    handleListLinks(w2, req2)
    resp2 := w2.Result()
    if resp2.StatusCode != http.StatusOK {
        t.Fatalf("expected 200 OK, got %d", resp2.StatusCode)
    }
    var got []Link
    if err := json.NewDecoder(resp2.Body).Decode(&got); err != nil {
        t.Fatal(err)
    }
    if len(got) != 1 || got[0].ShortName != "foo" || got[0].URL != "https://example.com" {
        t.Fatalf("unexpected list result: %#v", got)
    }
}

func TestCreateConflict(t *testing.T) {
    f := NewFakeStore()
    prev := store
    store = f
    defer func(){ store = prev }()

    l := Link{ShortName: "dup", URL: "https://a"}
    body, _ := json.Marshal(l)
    req := httptest.NewRequest("POST", "/api/links", bytes.NewReader(body))
    w := httptest.NewRecorder()
    handleCreateLink(w, req)
    if w.Result().StatusCode != http.StatusCreated { t.Fatalf("expected created") }

    // attempt duplicate
    req2 := httptest.NewRequest("POST", "/api/links", bytes.NewReader(body))
    w2 := httptest.NewRecorder()
    handleCreateLink(w2, req2)
    if w2.Result().StatusCode != http.StatusConflict { t.Fatalf("expected conflict") }
}
