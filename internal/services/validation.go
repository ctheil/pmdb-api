package services

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func ValidateStruct(data interface{}) (errs map[string]string) {
	if err := validate.Struct(data); err != nil {
		return ParseValidationErrors(err)
	}
	return nil
}

func ParseValidationErrors(validationAttempt error) (errs map[string]string) {
	if validationAttempt == nil {
		return nil
	}
	validationErrors, ok := validationAttempt.(validator.ValidationErrors)
	if !ok {
		return map[string]string{"error": "Invalid validation error"}
	}

	errs = make(map[string]string)

	for _, err := range validationErrors {
		field := err.Field()
		tag := err.Tag()
		errs[field] = "Validation failed on the '" + tag + "' tag"
	}
	return errs
}
