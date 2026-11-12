package httpx

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Code string

const (
	CodeBadRequest       Code = "bad_request"
	CodeValidationFailed Code = "validation_failed"
	CodeUnauthorized     Code = "unauthorized"
	CodeForbidden        Code = "forbidden"
	CodeNotFound         Code = "not_found"
	CodeConflict         Code = "conflict"
	CodeInternal         Code = "internal"
)

type ErrorBody struct {
	Code      Code              `json:"code"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields,omitempty"`
	RequestID string            `json:"requestId,omitempty"`
}

type Envelope struct {
	Error ErrorBody `json:"error"`
}

const ContextKeyRequestID = "requestID"

func Error(c *gin.Context, status int, body ErrorBody) {
	if id, exists := c.Get(ContextKeyRequestID); exists {
		value, ok := id.(string)
		if ok {
			body.RequestID = value
		}
	}
	c.AbortWithStatusJSON(status, Envelope{Error: body})
}

func WriteError(c *gin.Context, err error) {
	var ve *ValidationError
	switch {
	case errors.As(err, &ve):
		Error(c, http.StatusUnprocessableEntity, ErrorBody{
			Code:    CodeValidationFailed,
			Message: "invalid input",
			Fields:  ve.Fields,
		})
	case errors.Is(err, ErrUnauthorized):
		Error(c, http.StatusUnauthorized, ErrorBody{
			Code: CodeUnauthorized, Message: "authentication required",
		})
	case errors.Is(err, ErrForbidden):
		Error(c, http.StatusForbidden, ErrorBody{
			Code: CodeForbidden, Message: "forbidden access",
		})
	case errors.Is(err, ErrNotFound):
		Error(c, http.StatusNotFound, ErrorBody{
			Code: CodeNotFound, Message: "not found",
		})
	case errors.Is(err, ErrConflict):
		Error(c, http.StatusConflict, ErrorBody{Code: CodeConflict, Message: "conflict"})
	default:
		log.Printf("internal error: %v", err)
		Error(c, http.StatusInternalServerError, ErrorBody{Code: CodeInternal, Message: "internal server error"})
	}
}
