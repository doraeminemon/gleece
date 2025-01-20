package validation

import (
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/gopher-fleece/gleece/definitions"
)

// Global validator instance
var ValidatorInstance *validator.Validate

// Custom validation function to check if the slice is not nil
func validateNotNilSlice(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() == reflect.Slice {
		return !field.IsNil()
	}
	return false
}

// Custom validation function to check if a string starts with a letter
func validateStartsWithLetter(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if field == "" {
		return false
	}
	firstChar := rune(field[0])
	return unicode.IsLetter(firstChar)
}

// registerEnumValidator creates a custom validation function for enum types
func registerEnumValidator(enumValues interface{}) validator.Func {
	values := reflect.ValueOf(enumValues)
	allowedValues := make(map[interface{}]struct{})

	for i := 0; i < values.Len(); i++ {
		allowedValues[values.Index(i).Interface()] = struct{}{}
	}

	return func(fl validator.FieldLevel) bool {
		field := fl.Field().Interface()
		_, ok := allowedValues[field]
		return ok
	}
}

func validateRegex(fl validator.FieldLevel) bool {
	// Get the regex pattern from the tag's parameter
	params := strings.SplitN(fl.Param(), ":", 2)
	if len(params) < 1 {
		return false // Invalid tag usage
	}

	// Compile the regex
	re, err := regexp.Compile(params[0])
	if err != nil {
		return false // Invalid regex
	}

	// Validate the field value
	value := fl.Field().String()
	return re.MatchString(value)
}

func InitValidator() {
	// Initialize the validator instance
	ValidatorInstance = validator.New()

	// Register custom validation functions globally
	ValidatorInstance.RegisterValidation("not_nil_array", validateNotNilSlice)
	ValidatorInstance.RegisterValidation("starts_with_letter", validateStartsWithLetter)
	ValidatorInstance.RegisterValidation("regex", validateRegex)

	// Register enum validation functions

	// SecuritySchemeIn
	ValidatorInstance.RegisterValidation(
		"security_schema_in",
		registerEnumValidator(
			[]definitions.SecuritySchemeIn{
				definitions.InQuery,
				definitions.InHeader,
				definitions.InCookie,
			},
		),
	)

	// SecuritySchemeType
	ValidatorInstance.RegisterValidation(
		"security_schema_type",
		registerEnumValidator(
			[]definitions.SecuritySchemeType{
				definitions.APIKey,
				definitions.OAuth2,
				definitions.OpenIDConnect,
				definitions.HTTP,
			},
		),
	)
}
