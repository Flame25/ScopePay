package httpx

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

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

func jsonNameMap(dst any) map[string]string {
	out := map[string]string{}
	t := reflect.TypeOf(dst)
	for t != nil && t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t == nil || t.Kind() != reflect.Struct {
		return out
	}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		tag := f.Tag.Get("json")
		name := strings.Split(tag, ",")[0]
		if name == "" || name == "-" {
			name = f.Name
		}
		out[f.Name] = name
	}
	return out
}

func fieldErrors(dst any, err error) map[string]string {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return map[string]string{"body": "must be valid JSON matching the expected shape"}
	}

	names := jsonNameMap(dst)
	out := make(map[string]string, len(verrs))
	for _, fe := range verrs {
		key := fe.Field()
		if name, ok := names[key]; ok {
			key = name
		}
		out[key] = messageFor(fe)

	}
	return out
}

func BindJSON(c *gin.Context, dst any) error {
	if err := c.ShouldBindJSON(dst); err != nil {
		return &ValidationError{Fields: fieldErrors(dst, err)}
	}
	return nil
}
