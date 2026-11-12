package httpx

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Fields map[string]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

func messageFor(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fe.Param())
	case "oneof":
		return "must be one of: " + fe.Param()
	default:
		return "invalid value"
	}
}

func fieldErrors(dst any, err error) map[string]string {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return map[string]string{"body": "must be valid JSON matching the expected shape"}
	}
	out := make(map[string]string, len(verrs))
	for _, fe := range verrs {
		out[fe.Field()] = messageFor(fe)
	}
	return out
}

func BindJSON(c *gin.Context, dst any) error {
	if err := c.ShouldBindJSON(dst); err != nil {
		return &ValidationError{Fields: fieldErrors(dst, err)}
	}
	return nil
}
