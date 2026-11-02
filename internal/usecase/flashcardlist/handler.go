package flashcardlist

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	svc *service
}

func NewHandler(svc *service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	context := req.Context()
	flashcards, err := h.svc.db.ListFlashcards(context, 10)
	if err != nil {
		h.svc.logger.Error("db.ListFlashcards:", "err", err)
		w.WriteHeader(500)
		return
	}
	type response struct {
		Data []flashcard `json:"data"`
	}
	res, err := json.Marshal(response{Data: flashcards})
	if err != nil {
		h.svc.logger.Error("json.Marshal:", "err", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	_, _ = w.Write(res)
}
