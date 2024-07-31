package service

import (
	"encoding/json"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type validation struct {
	Validate   *validator.Validate
	Translator ut.Translator
}

var v *validation

func init() {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	validator := validator.New(validator.WithRequiredStructEnabled())
	en_translations.RegisterDefaultTranslations(validator, trans)

	v = &validation{
		Validate:   validator,
		Translator: trans,
	}
}

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	bb, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}
	return string(bb)
}

func (v *validation) Struct(a any) error {
	err := v.Validate.Struct(a)
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// This error can be either nil or InvalidValidationError.
		return err
	}

	var ve ValidationErrors
	for _, err := range errs {
		if v.Translator == nil {
			ve = append(ve, ValidationError{
				Field:   err.Field(),
				Message: err.Error(),
			})
		} else {
			ve = append(ve, ValidationError{
				Field:   err.Field(),
				Message: err.Translate(v.Translator),
			})
		}
	}

	return ve
}
