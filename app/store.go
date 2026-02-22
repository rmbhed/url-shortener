package main

import (
    "context"

    "cloud.google.com/go/firestore"
    "google.golang.org/api/iterator"
)

// Link represents a short link mapping
type Link struct {
    ShortName string `json:"shortName" firestore:"shortName"`
    URL       string `json:"url" firestore:"url"`
}

// LinkStore is an interface for storing and retrieving links.
type LinkStore interface {
    GetAll(ctx context.Context) ([]Link, error)
    Create(ctx context.Context, l Link) error
    Get(ctx context.Context, shortName string) (Link, error)
}

// FirestoreStore implements LinkStore using Firestore.
type FirestoreStore struct{
    Client *firestore.Client
}

func (f *FirestoreStore) GetAll(ctx context.Context) ([]Link, error) {
    iter := f.Client.Collection("urls").Documents(ctx)
    var links []Link
    for {
        doc, err := iter.Next()
        if err == iterator.Done { break }
        if err != nil { return nil, err }
        var l Link
        if err := doc.DataTo(&l); err != nil { return nil, err }
        links = append(links, l)
    }
    return links, nil
}

func (f *FirestoreStore) Create(ctx context.Context, l Link) error {
    _, err := f.Client.Collection("urls").Doc(l.ShortName).Create(ctx, l)
    return err
}

func (f *FirestoreStore) Get(ctx context.Context, shortName string) (Link, error) {
    doc, err := f.Client.Collection("urls").Doc(shortName).Get(ctx)
    if err != nil { return Link{}, err }
    var l Link
    if err := doc.DataTo(&l); err != nil { return Link{}, err }
    return l, nil
}
