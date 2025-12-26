package middleware

import (
	"net/http"

	"github.com/comply360/shared/errors"
	"github.com/comply360/shared/validator"
	"github.com/gin-gonic/gin"
)

// ValidateRequest is a middleware that validates request body against a struct
func ValidateRequest(v *validator.Validator, dataType interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind JSON to the data type
		if err := c.ShouldBindJSON(dataType); err != nil {
			c.JSON(http.StatusBadRequest, errors.BadRequest("Invalid JSON format: "+err.Error()))
			c.Abort()
			return
		}

		// Validate the bound data
		if apiErr := v.Validate(dataType); apiErr != nil {
			c.JSON(http.StatusBadRequest, apiErr)
			c.Abort()
			return
		}

		// Store validated data in context
		c.Set("validated_data", dataType)
		c.Next()
	}
}

// ValidateQuery validates query parameters
func ValidateQuery(v *validator.Validator, dataType interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Bind query parameters
		if err := c.ShouldBindQuery(dataType); err != nil {
			c.JSON(http.StatusBadRequest, errors.BadRequest("Invalid query parameters: "+err.Error()))
			c.Abort()
			return
		}

		// Validate
		if apiErr := v.Validate(dataType); apiErr != nil {
			c.JSON(http.StatusBadRequest, apiErr)
			c.Abort()
			return
		}

		c.Set("validated_query", dataType)
		c.Next()
	}
}
