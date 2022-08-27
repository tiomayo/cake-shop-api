package middlewares

import (
	"cake-store/internal/helpers"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func UseCustomValidatorHandler(e *echo.Echo) {
	newValidator := validator.New()
	e.Validator = &customValidator{validator: newValidator}

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if castedObject, ok := err.(validator.ValidationErrors); ok {
			var MessageValidation []helpers.ErrorObject
			for _, err := range castedObject {
				errObject := helpers.ErrorObject{}
				switch err.Tag() {
				case "required", "required_with", "required_without", "required_unless":
					errObject.Name = err.Field()
					errObject.Message = fmt.Sprintf("%s is required", err.Field())
				case "url", "numeric":
					errObject.Name = err.Field()
					errObject.Message = fmt.Sprintf("%s is not valid %s",
						err.Field(), err.Tag())
				default:
					errObject.Name = err.Field()
					errObject.Message = fmt.Sprintf("Validation error on field %s", err.Field())
				}

				if errObject.Name != "" {
					MessageValidation = append(MessageValidation, errObject)
				}
			}
			c.JSON(http.StatusUnprocessableEntity, helpers.JSONResponse{Message: "The given data was invalid.", Errors: MessageValidation})
		} else if castedObject, ok := err.(*echo.HTTPError); ok {
			log.Println(castedObject.Message)
			c.JSON(castedObject.Code, helpers.JSONResponse{Message: fmt.Sprintf("%v", castedObject.Message)})
		} else {
			c.Logger().Error(err)
			c.JSON(http.StatusInternalServerError, helpers.JSONResponse{Message: fmt.Sprintf("%v", err.Error())})
		}
	}
}
