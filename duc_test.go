package sso_test

import (
	"fmt"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"

	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"testing"
)

var (
	uni           *ut.UniversalTranslator
	ValidateTrans ut.Translator
)

func TestMy(t *testing.T) {
	zh2 := zh.New()
	uni = ut.New(zh2, zh2)

	ValidateTrans, _ = uni.GetTranslator("zh")

	v := validator.New()
	if err := zh_translations.RegisterDefaultTranslations(v, ValidateTrans); err != nil {
	}

	type User struct {
		Username string `validate:"required"`
		Tagline  string `validate:"required,lt=10"`
		Tagline2 string `validate:"required,gt=1"`
	}

	user := User{
		Username: "Joeybloggs",
		Tagline:  "This tagline is way too long.",
		Tagline2: "1",
	}
	err:=v.Struct(user)
	if err != nil {

		// translate all error at once
		errs := err.(validator.ValidationErrors)

		// returns a map with key = namespace & value = translated error
		// NOTICE: 2 errors are returned and you'll see something surprising
		// translations are i18n aware!!!!
		// eg. '10 characters' vs '1 character'
		fmt.Println(errs.Translate(ValidateTrans))
	}
}
