package handler

import "miso/internal/storage"

type Handler struct {
	storage storage.Storage
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}
