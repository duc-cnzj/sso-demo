package form

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

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
func ErrorsToMap(errors validator.ValidationErrors, ut ut.Translator) map[string]string {
	var m = map[string]string{}
	for _, e := range errors {
		m[e.Field()] = e.Translate(ut)
	}

	return m
}
