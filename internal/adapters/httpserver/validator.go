package httpserver

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

func NewValidator() (*validator.Validate, ut.Translator, error) {
	locale := en.New()
	uni := ut.New(locale, locale)
	translator, ok := uni.GetTranslator("en")
	if !ok {
		return nil, nil, errors.New("fail to find translator")
	}
	validate := validator.New()
	err := validate.RegisterTranslation(
		"required",
		translator,
		func(uni ut.Translator) error {
			return uni.Add("required", "{0} is required", true)
		},
		func(uni ut.Translator, fe validator.FieldError) string {
			t, _ := uni.T("required", fe.Field())
			return t
		},
	)
	if err != nil {
		return nil, nil, err
	}
	err = validate.RegisterTranslation(
		"min",
		translator,
		func(uni ut.Translator) error {
			return uni.Add("min", "{0} must be >= {1}, currently {2}", true)
		},
		func(uni ut.Translator, fe validator.FieldError) string {
			t, _ := uni.T("min", fe.Field(), fe.Param(), fmt.Sprint(fe.Value()))
			return t
		},
	)
	if err != nil {
		return nil, nil, err
	}
	err = validate.RegisterTranslation(
		"max",
		translator,
		func(uni ut.Translator) error {
			return uni.Add("max", "{0} must be <= {1}, currently {2}", true)
		},
		func(uni ut.Translator, fe validator.FieldError) string {
			t, _ := uni.T("max", fe.Field(), fe.Param(), fmt.Sprint(fe.Value()))
			return t
		},
	)
	if err != nil {
		return nil, nil, err
	}
	return validate, translator, nil
}

func ValidationErrorMessage(err validator.ValidationErrors, translator ut.Translator) string {
	translatedErrs := err.Translate(translator)
	keys := []string{}
	for key := range translatedErrs {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	veMessages := []string{}
	for _, key := range keys {
		if errMessage := translatedErrs[key]; errMessage != "" {
			veMessages = append(veMessages, errMessage)
		}
	}
	return strings.Join(veMessages, ", ")
}
