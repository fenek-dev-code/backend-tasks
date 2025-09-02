package handler

import "github.com/fenek-dev-code/backend-tasks/internal/storage"

type Handler struct {
	UserHandler *UserHandler
}

func NewHandler(storage *storage.Storage) *Handler {
	return &Handler{UserHandler: NewUserHandler(storage)}
}