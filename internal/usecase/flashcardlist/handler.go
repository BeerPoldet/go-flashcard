package flashcardlist

import (
	"context"
	"encoding/base64"
	"encoding/json"
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
	Perform(context.Context, queryOptions) (listFlashcardQueryResult, error)
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
	cur, err := newQueryCursor(queryParams.After, queryParams.Sort, queryParams.Order)
	if err != nil {
		h.logger.Debug("newCursor:", "err", err)
		switch {
		case errors.Is(err, errCursorWithSortOrOrder):
			httpserver.WriteJSON(
				http.StatusBadRequest,
				httpserver.ErrorResponse{
					Code:        "validation_error",
					Description: "after cannot be combined with sort or order",
				},
				w,
				h.logger,
			)
		case errors.Is(err, errInvalidCursor):
			httpserver.WriteJSON(
				http.StatusBadRequest,
				httpserver.ErrorResponse{
					Code:        "validation_error",
					Description: "after cursor is malformed",
				},
				w,
				h.logger,
			)
		default:
			httpserver.WriteInternalServerError(w, h.logger)
		}
		return
	}
	queryOptions := queryOptions{
		limit:  queryParams.Limit,
		cursor: cur,
	}
	result, err := h.uc.Perform(ctx, queryOptions)
	if err != nil {
		h.logger.Error("uc.Perform:", "err", err)
		httpserver.WriteInternalServerError(w, h.logger)
		return
	}
	resp, err := getResponse(result)
	if err != nil {
		h.logger.Error("getResponse:", "err", err)
		httpserver.WriteInternalServerError(w, h.logger)
		return
	}
	httpserver.WriteJSON(
		http.StatusOK,
		resp,
		w,
		h.logger,
	)
}

func getResponse(result listFlashcardQueryResult) (response, error) {
	flashcardListResponse := make([]flashcardResponse, 0, len(result.data))
	for _, f := range result.data {
		flashcardListResponse = append(
			flashcardListResponse,
			f.toFlashcardListResponse(),
		)
	}
	var afterCursor *string
	if result.after != nil {
		after, err := json.Marshal(result.after)
		if err != nil {
			return response{}, fmt.Errorf("%w: %w", errMarshalCursor, err)
		}
		afterCursor = new(base64.RawURLEncoding.EncodeToString(after))
	}

	return response{Data: flashcardListResponse, After: afterCursor}, nil
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
	Data  []flashcardResponse `json:"data"`
	After *string             `json:"after"`
}

type flashcardResponse struct {
	ID        int32     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type queryParams struct {
	Limit int32  `form:"limit" validate:"omitempty,min=1,max=20"`
	Sort  string `form:"sort"  validate:"omitempty,oneof=created_at"`
	Order string `form:"order" validate:"omitempty,oneof=asc desc"`
	After string `form:"after" validate:"omitempty"`
}

var (
	errCursorWithSortOrOrder = errors.New("cursor used with sort/order")
	errInvalidCursor         = errors.New("invalid cursor")
	errMarshalCursor         = errors.New("marshal cursor")
)

func newQueryCursor(after string, sort string, order string) (cursor, error) {
	cur := cursor{}
	if after != "" && (sort != "" || order != "") {
		return cur, errCursorWithSortOrOrder
	}
	if after != "" {
		afterByte, err := base64.RawURLEncoding.DecodeString(after)
		if err != nil {
			return cur, errInvalidCursor
		}
		if err := json.Unmarshal(afterByte, &cur); err != nil {
			return cur, errInvalidCursor
		}
		if err := cur.validate(); err != nil {
			return cur, err
		}
	}
	if sort != "" {
		cur.Sort = newSort(sort)
	}
	if order != "" {
		cur.Order = newOrder(order)
	}
	return cur, nil
}

func (c cursor) validate() error {
	if c.Order != ascOrder && c.Order != descOrder {
		return errInvalidCursor
	}
	switch c.Sort {
	case createdAtSort:
		if c.Position == nil || c.Position.CreatedAt == nil {
			return errInvalidCursor
		}
		return nil
	default:
		return errInvalidCursor
	}
}

func newSort(key string) sort {
	switch key {
	case "created_at":
		return createdAtSort
	default:
		return ""
	}
}

func newOrder(key string) order {
	switch key {
	case "asc":
		return ascOrder
	default:
		return descOrder
	}
}
