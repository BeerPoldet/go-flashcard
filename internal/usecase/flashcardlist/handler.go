package flashcardlist

import (
	"net/http"

	"github.com/beerpoldet/flash-card-go/internal/adapters/httpserver"
)

type Handler struct {
	uc     useCase
	logger logger
}

func NewHandler(uc useCase, logger logger) *Handler {
	return &Handler{uc: uc, logger: logger}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	context := req.Context()
	flashcards, err := h.uc.Perform(context, 10)
	if err != nil {
		httpserver.WriteJSON(
			http.StatusInternalServerError,
			httpserver.ErrorResponse{
				Code:        "internal_server_error",
				Description: "unexpected error occured",
			},
			w,
			h.logger,
		)
	} else {
		httpserver.WriteJSON(
			http.StatusOK,
			flashcardListResponse{Data: flashcards},
			w,
			h.logger,
		)
	}
}

type flashcardListResponse struct {
	Data []flashcard `json:"data"`
}
