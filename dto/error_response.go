package dto

type ErrorResponse struct {
    Status  int    `json:"status"`
    Error   string `json:"error"`
    Message string `json:"message"`
}