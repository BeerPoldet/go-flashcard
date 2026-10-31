package listflashcards

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Handler struct {
	svc service
}

func NewHandler(svc service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	flashcards := h.svc.ListFlashcards()
	type response struct {
		Data []string `json:"data"`
	}
	res, err := json.Marshal(response{Data: flashcards})
	if err != nil {
		w.WriteHeader(500)
		_, _ = fmt.Fprintf(w, "Something went wrong, error %v", err)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write(res)
}
