package flashcardlist

import (
	"net/http"
)

type Handler struct {
	uc useCaseInputPort
}

func NewHandler(uc useCaseInputPort) *Handler {
	return &Handler{uc: uc}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	context := req.Context()
	h.uc.listFlashcards(context, 10)
}
