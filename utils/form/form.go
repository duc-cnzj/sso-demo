package form

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type ValidateError struct {
	Field string
	Msg string
}

type ValidateErrors []ValidateError

//v.Value():
//v.ActualTag():  required
//v.Field():  Password
//v.Namespace():  LoginForm.Password
//v.Type():  string
//v.StructField():  Password
//v.Tag():  required
//v.Kind():  string
//v.Param():
//v.StructNamespace():  LoginForm.Password
//v.Value():
//v.ActualTag():  required
//v.Field():  RedirectUrl
//v.Namespace():  LoginForm.RedirectUrl
//v.Type():  string
//v.StructField():  RedirectUrl
//v.Tag():  required
//v.Kind():  string
//v.Param():
//v.StructNamespace():  LoginForm.RedirectUrl
func ErrorsToMap(errors interface{}, ut ut.Translator) map[string]string {
	switch errors.(type) {
	case validator.ValidationErrors:
		var m = map[string]string{}
		for _, e := range errors.(validator.ValidationErrors) {
			m[e.Field()] = e.Translate(ut)
		}

		return m
	case ValidateErrors:
		var m = map[string]string{}
		for _, e := range errors.(ValidateErrors) {
			m[e.Field] = e.Msg
		}

		return m
	default:
		panic(fmt.Sprintf("not support %T", errors))
	}
}
