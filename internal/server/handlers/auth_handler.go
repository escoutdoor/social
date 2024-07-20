package handlers

import "github.com/escoutdoor/social/internal/db/store"

type AuthHandler struct {
	store store.AuthStorer
}

func NewAuthHandler(store store.AuthStorer) AuthHandler {
	return AuthHandler{
		store: store,
	}
}
