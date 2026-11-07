package flashcardlist

import (
	"context"
	"net/http"
	"time"

	"github.com/beerpoldet/flash-card-go/internal/adapters/httpserver"
)

type usecase interface {
	Perform(context.Context, int) ([]flashcard, error)
}

type Handler struct {
	uc     usecase
	logger logger
}

func newHandler(uc usecase, logger logger) *Handler {
	return &Handler{uc: uc, logger: logger}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	flashcards, err := h.uc.Perform(ctx, 10)
	if err != nil {
		h.logger.Error("uc.Perform:", "err", err)
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
		flashcardListResponse := make([]flashcardResponse, 0, len(flashcards))
		for _, f := range flashcards {
			flashcardListResponse = append(
				flashcardListResponse,
				f.toFlashcardListResponse(),
			)
		}
		httpserver.WriteJSON(
			http.StatusOK,
			response{Data: flashcardListResponse},
			w,
			h.logger,
		)
	}
}

func (f *flashcard) toFlashcardListResponse() flashcardResponse {
	return flashcardResponse{
		ID:        f.ID,
		Title:     f.Title,
		Content:   f.Content,
		CreatedAt: f.CreatedAt,
	}
}

type response struct {
	Data []flashcardResponse `json:"data"`
}

type flashcardResponse struct {
	ID        int32     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
