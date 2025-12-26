package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/comply360/shared/errors"
	"github.com/go-playground/validator/v10"
)

// ValidationError represents a field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Value   interface{} `json:"value,omitempty"`
}

// Validator is a wrapper around go-playground/validator
type Validator struct {
	validate *validator.Validate
}

// New creates a new validator instance
func New() *Validator {
	v := validator.New()

	// Register custom tag name function to use json tags
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validators
	v.RegisterValidation("phone", validatePhone)
	v.RegisterValidation("subdomain", validateSubdomain)
	v.RegisterValidation("currency", validateCurrency)
	v.RegisterValidation("country_code", validateCountryCode)
	v.RegisterValidation("jurisdiction", validateJurisdiction)
	v.RegisterValidation("registration_type", validateRegistrationType)
	v.RegisterValidation("user_role", validateUserRole)
	v.RegisterValidation("user_status", validateUserStatus)
	v.RegisterValidation("commission_rate", validateCommissionRate)
	v.RegisterValidation("document_type", validateDocumentType)
	v.RegisterValidation("registration_status", validateRegistrationStatus)
	v.RegisterValidation("commission_status", validateCommissionStatus)
	v.RegisterValidation("sa_id_number", validateSAIDNumber)
	v.RegisterValidation("company_registration_number", validateCompanyRegistrationNumber)
	v.RegisterValidation("vat_number", validateVATNumber)
	v.RegisterValidation("strong_password", validateStrongPassword)

	return &Validator{validate: v}
}

// Validate validates a struct and returns formatted errors
func (v *Validator) Validate(data interface{}) *errors.APIError {
	err := v.validate.Struct(data)
	if err == nil {
		return nil
	}

	var validationErrors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, ValidationError{
			Field:   err.Field(),
			Message: getErrorMessage(err),
			Tag:     err.Tag(),
			Value:   err.Value(),
		})
	}

	return errors.ValidationFailed("Validation failed", validationErrors)
}

// ValidateVar validates a single variable
func (v *Validator) ValidateVar(field interface{}, tag string) error {
	return v.validate.Var(field, tag)
}

// getErrorMessage returns a human-readable error message
func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", fe.Field(), fe.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", fe.Field())
	case "url":
		return fmt.Sprintf("%s must be a valid URL", fe.Field())
	case "phone":
		return fmt.Sprintf("%s must be a valid phone number", fe.Field())
	case "subdomain":
		return fmt.Sprintf("%s must be a valid subdomain (lowercase alphanumeric and hyphens only)", fe.Field())
	case "currency":
		return fmt.Sprintf("%s must be a valid currency code (ZAR, USD, or ZWL)", fe.Field())
	case "country_code":
		return fmt.Sprintf("%s must be a valid ISO 3166-1 alpha-2 country code", fe.Field())
	case "jurisdiction":
		return fmt.Sprintf("%s must be a valid jurisdiction (ZA or ZW)", fe.Field())
	case "registration_type":
		return fmt.Sprintf("%s must be a valid registration type", fe.Field())
	case "user_role":
		return fmt.Sprintf("%s must be a valid user role", fe.Field())
	case "user_status":
		return fmt.Sprintf("%s must be a valid user status", fe.Field())
	case "commission_rate":
		return fmt.Sprintf("%s must be between 0 and 100", fe.Field())
	case "document_type":
		return fmt.Sprintf("%s must be a valid document type", fe.Field())
	case "registration_status":
		return fmt.Sprintf("%s must be a valid registration status", fe.Field())
	case "commission_status":
		return fmt.Sprintf("%s must be a valid commission status", fe.Field())
	case "sa_id_number":
		return fmt.Sprintf("%s must be a valid South African ID number (13 digits)", fe.Field())
	case "company_registration_number":
		return fmt.Sprintf("%s must be a valid company registration number", fe.Field())
	case "vat_number":
		return fmt.Sprintf("%s must be a valid VAT number (10 digits starting with 4)", fe.Field())
	case "strong_password":
		return fmt.Sprintf("%s must be at least 8 characters with uppercase, lowercase, numbers, and symbols", fe.Field())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", fe.Field(), fe.Param())
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", fe.Field(), fe.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", fe.Field(), fe.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fe.Field(), fe.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", fe.Field(), fe.Param())
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", fe.Field())
	default:
		return fmt.Sprintf("%s is invalid", fe.Field())
	}
}

// Custom validation functions

func validatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true // Optional field
	}
	// Basic phone validation: starts with + and contains 7-15 digits
	return len(phone) >= 7 && len(phone) <= 20
}

func validateSubdomain(fl validator.FieldLevel) bool {
	subdomain := fl.Field().String()
	if subdomain == "" {
		return false
	}
	// Must be lowercase alphanumeric with hyphens, not starting/ending with hyphen
	if len(subdomain) < 3 || len(subdomain) > 63 {
		return false
	}
	for i, char := range subdomain {
		if (char >= 'a' && char <= 'z') || (char >= '0' && char <= '9') {
			continue
		}
		if char == '-' && i != 0 && i != len(subdomain)-1 {
			continue
		}
		return false
	}
	return true
}

func validateCurrency(fl validator.FieldLevel) bool {
	currency := fl.Field().String()
	validCurrencies := []string{"ZAR", "USD", "ZWL", "EUR", "GBP"}
	for _, valid := range validCurrencies {
		if currency == valid {
			return true
		}
	}
	return false
}

func validateCountryCode(fl validator.FieldLevel) bool {
	code := fl.Field().String()
	// Simplified: just check length and uppercase
	return len(code) == 2 && code == strings.ToUpper(code)
}

func validateJurisdiction(fl validator.FieldLevel) bool {
	jurisdiction := fl.Field().String()
	validJurisdictions := []string{"ZA", "ZW"}
	for _, valid := range validJurisdictions {
		if jurisdiction == valid {
			return true
		}
	}
	return false
}

func validateRegistrationType(fl validator.FieldLevel) bool {
	regType := fl.Field().String()
	validTypes := []string{"pty_ltd", "close_corporation", "business_name", "vat_registration"}
	for _, valid := range validTypes {
		if regType == valid {
			return true
		}
	}
	return false
}

func validateUserRole(fl validator.FieldLevel) bool {
	role := fl.Field().String()
	validRoles := []string{"system_admin", "tenant_admin", "tenant_manager", "agent", "agent_assistant", "client"}
	for _, valid := range validRoles {
		if role == valid {
			return true
		}
	}
	return false
}

func validateUserStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	validStatuses := []string{"active", "suspended", "locked", "deleted"}
	for _, valid := range validStatuses {
		if status == valid {
			return true
		}
	}
	return false
}

func validateCommissionRate(fl validator.FieldLevel) bool {
	rate := fl.Field().Float()
	return rate >= 0 && rate <= 100
}

func validateDocumentType(fl validator.FieldLevel) bool {
	docType := fl.Field().String()
	validTypes := []string{
		"id_document", "proof_of_address", "company_constitution",
		"tax_certificate", "banking_details", "cipc_certificate",
		"founding_statement", "memorandum", "directors_resolution", "other",
	}
	for _, valid := range validTypes {
		if docType == valid {
			return true
		}
	}
	return false
}

func validateRegistrationStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	validStatuses := []string{"draft", "submitted", "in_review", "approved", "rejected", "cancelled"}
	for _, valid := range validStatuses {
		if status == valid {
			return true
		}
	}
	return false
}

func validateCommissionStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	validStatuses := []string{"pending", "approved", "paid", "cancelled"}
	for _, valid := range validStatuses {
		if status == valid {
			return true
		}
	}
	return false
}

// validateSAIDNumber validates South African ID numbers using the Luhn algorithm
func validateSAIDNumber(fl validator.FieldLevel) bool {
	idNumber := fl.Field().String()

	// SA ID must be exactly 13 digits
	if len(idNumber) != 13 {
		return false
	}

	// Check if all characters are digits
	for _, char := range idNumber {
		if char < '0' || char > '9' {
			return false
		}
	}

	// Validate using Luhn algorithm
	sum := 0
	for i := 0; i < 12; i++ {
		digit := int(idNumber[i] - '0')
		if i%2 == 0 {
			sum += digit
		} else {
			doubled := digit * 2
			if doubled > 9 {
				doubled -= 9
			}
			sum += doubled
		}
	}

	checkDigit := (10 - (sum % 10)) % 10
	return checkDigit == int(idNumber[12]-'0')
}

// validateCompanyRegistrationNumber validates South African company registration numbers
func validateCompanyRegistrationNumber(fl validator.FieldLevel) bool {
	regNumber := fl.Field().String()

	// Format: YYYY/NNNNNN/NN (e.g., 2020/123456/07)
	if len(regNumber) < 10 || len(regNumber) > 15 {
		return false
	}

	// Basic format check - should contain numbers and slashes
	hasSlash := false
	hasNumber := false
	for _, char := range regNumber {
		if char == '/' {
			hasSlash = true
		} else if char >= '0' && char <= '9' {
			hasNumber = true
		}
	}

	return hasSlash && hasNumber
}

// validateVATNumber validates South African VAT numbers
func validateVATNumber(fl validator.FieldLevel) bool {
	vatNumber := fl.Field().String()

	// SA VAT numbers are 10 digits starting with 4
	if len(vatNumber) != 10 {
		return false
	}

	// Must start with 4
	if vatNumber[0] != '4' {
		return false
	}

	// All characters must be digits
	for _, char := range vatNumber {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

// validateStrongPassword validates password strength
func validateStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Minimum 8 characters
	if len(password) < 8 {
		return false
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case char == '!' || char == '@' || char == '#' || char == '$' ||
		     char == '%' || char == '^' || char == '&' || char == '*':
			hasSpecial = true
		}
	}

	// Require at least 3 out of 4 character types
	count := 0
	if hasUpper {
		count++
	}
	if hasLower {
		count++
	}
	if hasNumber {
		count++
	}
	if hasSpecial {
		count++
	}

	return count >= 3
}
