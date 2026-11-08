package flashcardlist

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/beerpoldet/flash-card-go/internal/adapters/httpserver"
	"github.com/go-playground/form/v4"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type usecase interface {
	Perform(context.Context, int) ([]flashcard, error)
}

type Handler struct {
	uc         usecase
	logger     logger
	decoder    *form.Decoder
	validate   *validator.Validate
	translator ut.Translator
}

func newHandler(uc usecase, logger logger) *Handler {
	validate, translator, err := httpserver.NewValidator()
	if err != nil {
		panic(err)
	}
	return &Handler{
		uc:         uc,
		logger:     logger,
		decoder:    form.NewDecoder(),
		validate:   validate,
		translator: translator,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	queryParams := queryParams{
		Limit: 10,
	}
	if err := h.decoder.Decode(&queryParams, req.URL.Query()); err != nil {
		h.logger.Debug("decoder queryParams:", "err", err)
		httpserver.WriteJSON(
			http.StatusBadRequest,
			httpserver.ErrorResponse{
				Code:        "decode_error",
				Description: "query params decoder error",
			},
			w,
			h.logger,
		)
		return
	}
	if err := h.validate.Struct(&queryParams); err != nil {
		ves, ok := errors.AsType[validator.ValidationErrors](err)
		if !ok {
			h.logger.Error("validate.Struct:", "err", err)
			httpserver.WriteInternalServerError(w, h.logger)
			return
		}
		h.logger.Debug("validation queryParams:", "err", err)
		httpserver.WriteJSON(
			http.StatusBadRequest,
			httpserver.ErrorResponse{
				Code: "validation_error",
				Description: fmt.Sprintf(
					"query params validation error: %s",
					httpserver.ValidationErrorMessage(ves, h.translator),
				),
			},
			w,
			h.logger,
		)
		return
	}
	flashcards, err := h.uc.Perform(ctx, queryParams.Limit)
	if err != nil {
		h.logger.Error("uc.Perform:", "err", err)
		httpserver.WriteInternalServerError(w, h.logger)
		return
	}
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

type queryParams struct {
	Limit int    `form:"limit" validate:"min=1,max=20"`
	Next  string `form:"next"`
}
