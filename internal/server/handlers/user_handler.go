package handlers

import "github.com/escoutdoor/social/internal/db/store"

type UserHandler struct {
	store store.UserStorer
}

func NewUserHandler(store store.UserStorer) UserHandler {
	return UserHandler{
		store: store,
	}
}
