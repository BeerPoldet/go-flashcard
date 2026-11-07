package httpserver

type ErrorResponse struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}
