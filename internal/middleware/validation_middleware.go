package middleware

import (
	"errors"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// ValidatorMiddleware handles request validation
func ValidatorMiddleware(model interface{}) fiber.Handler {
	validate := validator.New()

	return func(c *fiber.Ctx) error {
		// Skip validation for GET and DELETE requests with no body
		if c.Method() == fiber.MethodGet || c.Method() == fiber.MethodDelete {
			return c.Next()
		}

		// Create a new instance of the model
		input := reflect.New(reflect.TypeOf(model)).Elem()
		
		// Validate body
		if err := c.BodyParser(input.Addr().Interface()); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Validate struct
		if err := validate.Struct(input.Interface()); err != nil {
			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				var errorMessages []string
				for _, e := range validationErrors {
					errorMessages = append(errorMessages, formatValidationError(e))
				}

				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"errors": errorMessages,
				})
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Store validated input in context for route handlers
		c.Locals("validatedInput", input.Interface())
		return c.Next()
	}
}

// formatValidationError creates a human-readable validation error
func formatValidationError(e validator.FieldError) string {
	field := e.Field()
	tag := e.Tag()
	param := e.Param()

	switch tag {
	case "required":
		return field + " is required"
	case "email":
		return field + " must be a valid email address"
	case "min":
		return field + " must be at least " + param
	case "max":
		return field + " must be at most " + param
	default:
		return field + " failed " + tag + " validation"
	}
}
