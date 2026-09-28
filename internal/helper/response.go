package helper

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strings"

)

// Meta adalah informasi tentang response (code, status, pesan).
// Selalu ada di setiap response, baik sukses maupun error.
type Meta struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}


type APIResponse struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

func Success(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, APIResponse{
		Data: data,
		Meta: Meta{
			Code:    code,
			Status:  "success",
			Message: message,
		},
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, APIResponse{
		Data: nil,
		Meta: Meta{
			Code:    code,
			Status:  "error",
			Message: message,
		},
	})
}

func FormatValidatorError(err error) string {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errorMessages []string
		for _, e := range validationErrors {
			switch e.Tag() {
			case "required":
				errorMessages = append(errorMessages, e.Field()+" tidak boleh kosong")
			case "email":
				errorMessages = append(errorMessages, e.Field()+" harus berupa format email yang valid")
			case "min":
				errorMessages = append(errorMessages, e.Field()+" terlalu pendek")
			default:
				errorMessages = append(errorMessages, e.Field()+" tidak valid")
			}
		}
		return strings.Join(errorMessages, ", ")
	}
	return err.Error()
}
