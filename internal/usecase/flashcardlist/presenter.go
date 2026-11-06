package flashcardlist

import (
	"encoding/json"
	"net/http"
)

type flashcardListResponse struct {
	Data  []flashcard `json:"data"`
	Error string      `json:"error"`
}

type presenter struct {
	logger logger
	writer http.ResponseWriter
}

func newPresenter(logger logger) *presenter {
	return &presenter{logger: logger}
}

func (p *presenter) getFlashcardListResponse(data []flashcard, err error) {
	res, err := json.Marshal(toFlashcardListResponse(data, err))
	if err != nil {
		p.logger.Error("json.Marshal:", "err", err)
		p.writer.WriteHeader(500)
		return
	}
	p.writer.WriteHeader(200)
	_, _ = p.writer.Write(res)
}

func toFlashcardListResponse(data []flashcard, err error) flashcardListResponse {
	if err != nil {
		return flashcardListResponse{
			Error: error.Error(err),
		}
	}
	return flashcardListResponse{
		Data: data,
	}
}
